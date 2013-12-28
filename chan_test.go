package shellac

import (
	"io"
	"testing"
)

func TestChanReader(t *testing.T) {
	ch := make(chan string)
	r := NewChanReader(ch)
	go func() {
		for i := 0; i < 3; i++ {
			ch <- "hi"
		}
		close(ch)
	}()
	p := make([]byte, 4)
	for i := 0; i < 3; i++ {
		n, err := r.Read(p)
		if nil != err {
			t.Fatal(err)
		}
		if 3 != n {
			t.Fatal(n)
		}
		if "hi\n" != string(p[:n]) {
			t.Fatal(p)
		}
	}
	if _, err := r.Read(p); io.EOF != err {
		t.Fatal(err)
	}
}

func TestChanReaderOverflow(t *testing.T) {
	ch := make(chan string)
	r := NewChanReader(ch)
	go func() {
		ch <- "hi"
		close(ch)
	}()
	p := make([]byte, 1)
	for _, s := range []string{"h", "i", "\n"} {
		n, err := r.Read(p)
		if nil != err {
			t.Fatal(err)
		}
		if len(s) != n {
			t.Fatal(n)
		}
		if s != string(p[:n]) {
			t.Fatal(p)
		}
	}
	if _, err := r.Read(p); io.EOF != err {
		t.Fatal(err)
	}
}

func TestChanWriter(t *testing.T) {
	ch := make(chan string)
	w := NewChanWriter(ch)
	go func() {
		for i := 0; i < 3; i++ {
			if _, err := w.Write([]byte("hi\n")); nil != err {
				t.Fatal(err)
			}
		}
		w.Close()
	}()
	for i := 0; i < 3; i++ {
		s, ok := <-ch
		if !ok {
			t.Fatal(ok)
		}
		if "hi" != s {
			t.Fatal(s)
		}
	}
	if _, ok := <-ch; ok {
		t.Fatal(ok)
	}
}

func TestChanWriterOverflow(t *testing.T) {
	ch := make(chan string)
	w := NewChanWriter(ch)
	go func() {
		if _, err := w.Write([]byte("hi\nhi\nhi\n")); nil != err {
			t.Fatal(err)
		}
		w.Close()
	}()
	for i := 0; i < 3; i++ {
		s, ok := <-ch
		if !ok {
			t.Fatal(ok)
		}
		if "hi" != s {
			t.Fatal(s)
		}
	}
	if _, ok := <-ch; ok {
		t.Fatal(ok)
	}
}

func TestChanWriterUnderflow(t *testing.T) {
	ch := make(chan string)
	w := NewChanWriter(ch)
	go func() {
		for _, b := range []byte{'h', 'i', '\n'} {
			if _, err := w.Write([]byte{b}); nil != err {
				t.Fatal(err)
			}
		}
		w.Close()
	}()
	s, ok := <-ch
	if !ok {
		t.Fatal(ok)
	}
	if "hi" != s {
		t.Fatal(s)
	}
	if _, ok := <-ch; ok {
		t.Fatal(ok)
	}
}
