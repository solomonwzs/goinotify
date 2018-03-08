package goinotify

import (
	"sync"
	"syscall"
)

type Watcher struct {
	inotifyFd int
	epfd      int

	end     chan struct{}
	endLock *sync.Mutex
}

func NewWatcher(flags int) (w *Watcher, err error) {
	w = &Watcher{
		end:     make(chan struct{}),
		endLock: &sync.Mutex{},
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
	go w.serv()

	return
}

func (w *Watcher) serv() {
	events := make([]syscall.EpollEvent, _MAX_EVENTS, _MAX_EVENTS)
	buffer := make([]byte, _MAX_BUFFER_SIZE, _MAX_BUFFER_SIZE)
	for {
		nevents, err := syscall.EpollWait(w.epfd, events, _EPOLL_WAIT_TIMEOUT_MS)
		if err != nil {
			return
		}

		for i := 0; i < nevents; i++ {
			n, err := syscall.Read(int(events[i].Fd), buffer)
			if err != nil {
				return
			}
			// var e syscall.InotifyEvent
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
