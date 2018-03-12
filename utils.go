package goinotify

import (
	"errors"
	"syscall"
)

const (
	_MAX_EVENTS            = 100
	_MAX_BUFFER_SIZE       = 0xffff
	_EPOLL_WAIT_TIMEOUT_MS = 500

	IN_ACCESS        = syscall.IN_ACCESS
	IN_ALL_EVENTS    = syscall.IN_ALL_EVENTS
	IN_ATTRIB        = syscall.IN_ATTRIB
	IN_CLASSA_HOST   = syscall.IN_CLASSA_HOST
	IN_CLASSA_MAX    = syscall.IN_CLASSA_MAX
	IN_CLASSA_NET    = syscall.IN_CLASSA_NET
	IN_CLASSA_NSHIFT = syscall.IN_CLASSA_NSHIFT
	IN_CLASSB_HOST   = syscall.IN_CLASSB_HOST
	IN_CLASSB_MAX    = syscall.IN_CLASSB_MAX
	IN_CLASSB_NET    = syscall.IN_CLASSB_NET
	IN_CLASSB_NSHIFT = syscall.IN_CLASSB_NSHIFT
	IN_CLASSC_HOST   = syscall.IN_CLASSC_HOST
	IN_CLASSC_NET    = syscall.IN_CLASSC_NET
	IN_CLASSC_NSHIFT = syscall.IN_CLASSC_NSHIFT
	IN_CLOEXEC       = syscall.IN_CLOEXEC
	IN_CLOSE         = syscall.IN_CLOSE
	IN_CLOSE_NOWRITE = syscall.IN_CLOSE_NOWRITE
	IN_CLOSE_WRITE   = syscall.IN_CLOSE_WRITE
	IN_CREATE        = syscall.IN_CREATE
	IN_DELETE        = syscall.IN_DELETE
	IN_DELETE_SELF   = syscall.IN_DELETE_SELF
	IN_DONT_FOLLOW   = syscall.IN_DONT_FOLLOW
	IN_EXCL_UNLINK   = syscall.IN_EXCL_UNLINK
	IN_IGNORED       = syscall.IN_IGNORED
	IN_ISDIR         = syscall.IN_ISDIR
	IN_LOOPBACKNET   = syscall.IN_LOOPBACKNET
	IN_MASK_ADD      = syscall.IN_MASK_ADD
	IN_MODIFY        = syscall.IN_MODIFY
	IN_MOVE          = syscall.IN_MOVE
	IN_MOVED_FROM    = syscall.IN_MOVED_FROM
	IN_MOVED_TO      = syscall.IN_MOVED_TO
	IN_MOVE_SELF     = syscall.IN_MOVE_SELF
	IN_NONBLOCK      = syscall.IN_NONBLOCK
	IN_ONESHOT       = syscall.IN_ONESHOT
	IN_ONLYDIR       = syscall.IN_ONLYDIR
	IN_OPEN          = syscall.IN_OPEN
	IN_Q_OVERFLOW    = syscall.IN_Q_OVERFLOW
	IN_UNMOUNT       = syscall.IN_UNMOUNT
)

var (
	ERR_WATCHER_WAS_CLOSED = errors.New("goinotify: watcher was closed")
	ERR_TIMEOUT            = errors.New("goinotify: timeout")
)
