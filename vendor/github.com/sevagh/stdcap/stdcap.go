package stdcap

import (
	"bytes"
	"io"
	"os"
	"sync"
)

type stdcap struct {
	out bool
	mu  sync.RWMutex
}

var (
	sOut    *stdcap
	sErr    *stdcap
	onceOut sync.Once
	onceErr sync.Once
)

// StdoutCapture initializes an stdcap object for stdout
func StdoutCapture() *stdcap {
	onceOut.Do(func() {
		sOut = &stdcap{
			out: true,
		}
	})
	return sOut
}

// StderrCapture initializes an stdcap object for stderr
func StderrCapture() *stdcap {
	onceErr.Do(func() {
		sErr = &stdcap{
			out: false,
		}
	})
	return sErr
}

// Capture executes f() and returns the captured data
func (s *stdcap) Capture(f func()) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var old, r, w *os.File

	if s.out {
		old = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w
	} else {
		old = os.Stderr
		r, w, _ = os.Pipe()
		os.Stderr = w
	}

	f()

	outC := make(chan string)
	defer close(outC)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()

	if s.out {
		os.Stdout = old
	} else {
		os.Stderr = old
	}

	return <-outC
}
