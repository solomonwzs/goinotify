package goinotify

/*
#include <sys/inotify.h>
*/
import "C"
import (
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type Watcher struct {
	inotifyFd int
	epfd      int

	end     chan struct{}
	endLock *sync.Mutex

	events      []InotifyEventRaw
	eventsLock  *sync.Mutex
	eventNotify chan struct{}
}

func NewWatcher(flags int) (w *Watcher, err error) {
	w = &Watcher{
		end:     make(chan struct{}),
		endLock: &sync.Mutex{},

		events:      []InotifyEventRaw{},
		eventsLock:  &sync.Mutex{},
		eventNotify: make(chan struct{}),
	}

	if w.inotifyFd, err = syscall.InotifyInit1(flags); err != nil {
		return nil, err
	}

	if w.epfd, err = syscall.EpollCreate(1); err != nil {
		syscall.Close(w.inotifyFd)
		return nil, err
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(w.inotifyFd)
	if err = syscall.EpollCtl(w.epfd, syscall.EPOLL_CTL_ADD, w.inotifyFd,
		&event); err != nil {
		syscall.Close(w.inotifyFd)
		syscall.Close(w.epfd)
		return
	}
	go w.readEvents()

	return
}

func (w *Watcher) readEvents() {
	events := make([]syscall.EpollEvent, _MAX_EVENTS, _MAX_EVENTS)
	buffer := make([]byte, _MAX_BUFFER_SIZE, _MAX_BUFFER_SIZE)
	offset := 0
	for !w.IsClose() {
		nevents, err := syscall.EpollWait(w.epfd, events,
			_EPOLL_WAIT_TIMEOUT_MS)
		if err != nil {
			return
		}

		for i := 0; i < nevents; i++ {
			if events[i].Events&syscall.EPOLLIN != 0 {
				n, err := syscall.Read(int(events[i].Fd), buffer[offset:])
				if err != nil {
					return
				}
				n += offset

				w.eventsLock.Lock()
				j := 0
				for j < n {
					if n-j < syscall.SizeofInotifyEvent {
						offset = copy(buffer, buffer[j:])
						break
					}

					e := (*syscall.InotifyEvent)(unsafe.Pointer(&buffer[j]))
					size := syscall.SizeofInotifyEvent + e.Len
					if n-j < int(size) {
						offset = copy(buffer, buffer[j:])
						break
					}

					raw := InotifyEventRaw(make([]byte, size, size))
					copy(raw, buffer[j:])
					w.events = append(w.events, raw)
					j += int(size)
				}
				if j == n {
					offset = 0
				}
				w.notifyEvents()
				w.eventsLock.Unlock()
			}
		}
	}
}

func (w *Watcher) notifyEvents() {
	select {
	case w.eventNotify <- struct{}{}:
	default:
	}
}

func (w *Watcher) GetEvent(timeout time.Duration) (r InotifyEventRaw, err error) {
	for {
		r = nil
		w.eventsLock.Lock()
		if len(w.events) > 0 {
			r = w.events[0]
			w.events = w.events[1:]

			if len(w.events) > 0 {
				w.notifyEvents()
			}
		}
		w.eventsLock.Unlock()

		if r != nil {
			return
		}

		var deadline <-chan time.Time
		if timeout > 0 {
			timer := time.NewTimer(timeout)
			defer timer.Stop()
			deadline = timer.C

		}
		select {
		case <-w.eventNotify:
			break
		case <-deadline:
			return nil, ERR_TIMEOUT
		case <-w.end:
			return nil, ERR_WATCHER_WAS_CLOSED
		}
	}
}

func (w *Watcher) AddWatch(pathname string, mask uint32) (watchdesc int, err error) {
	return syscall.InotifyAddWatch(w.inotifyFd, pathname, mask)
}

func (w *Watcher) DelWatch(watchdesc int) (err error) {
	_, err = syscall.InotifyRmWatch(w.inotifyFd, uint32(watchdesc))
	return
}

func (w *Watcher) Close() error {
	w.endLock.Lock()
	defer w.endLock.Unlock()

	if w.IsClose() {
		return ERR_WATCHER_WAS_CLOSED
	}
	close(w.end)
	syscall.Close(w.inotifyFd)
	syscall.Close(w.epfd)
	return nil
}

func (w *Watcher) IsClose() bool {
	select {
	case <-w.end:
		return true
	default:
		return false
	}
}
