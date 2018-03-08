package goinotify

import (
	"errors"
	"syscall"
)

const (
	_MAX_EVENTS            = 100
	_MAX_BUFFER_SIZE       = 0xffff
	_EPOLL_WAIT_TIMEOUT_MS = 500

	IN_NONBLOCK = syscall.IN_NONBLOCK
	IN_CLOEXEC  = syscall.IN_CLOEXEC
)

var (
	ERR_WATCHER_WAS_CLOSED = errors.New("goinotify: watcher was closed")
)
