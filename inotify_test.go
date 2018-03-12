package goinotify

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestInotify(t *testing.T) {
	w, err := NewWatcher(0)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			fd, err := os.OpenFile("/tmp/1.txt", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				continue
			}

			buf := make([]byte, 16)
			fd.Read(buf)
			fd.WriteString("hello")
			fd.Close()

			os.Remove("/tmp/1.txt")
		}
	}()

	w.AddWatch("/tmp", syscall.IN_ALL_EVENTS)
	for {
		if r, err := w.GetEvent(2 * time.Second); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("file: %s, mask: %x\n", r.Name(), r.Mask())
		}
	}
	wg.Wait()
}
