package goinotify

import (
	"syscall"
	"unsafe"
)

type InotifyEventRaw []byte

func (r InotifyEventRaw) Wd() int32 {
	return (*syscall.InotifyEvent)(unsafe.Pointer(&r[0])).Wd
}

func (r InotifyEventRaw) Mask() uint32 {
	return (*syscall.InotifyEvent)(unsafe.Pointer(&r[0])).Mask
}

func (r InotifyEventRaw) Cookie() uint32 {
	return (*syscall.InotifyEvent)(unsafe.Pointer(&r[0])).Cookie
}

func (r InotifyEventRaw) Name() string {
	return string(r[syscall.SizeofInotifyEvent:])
}
