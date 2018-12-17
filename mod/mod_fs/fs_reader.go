package mod_fs

import (
	"os"
	"bufio"
	"fmt"
)

type reader struct {
	f *os.File
	r *bufio.Reader
}

func createReaderModule(path string) (*reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r := &reader{f:f, r:bufio.NewReader(f)}
	return r, nil
}

func (m *reader) Available() int {
	return m.r.Buffered()
}

func (m *reader) Skip(n int) (int, error) {
	return m.r.Discard(n)
}

func (m *reader) ReadLine() ([]byte, error) {
	return m.r.ReadBytes('\n')
}

func (m *reader) ReadStringLine() (string, error) {
	return m.r.ReadString('\n')
}

func (m *reader) Read(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be > 0")
	}
	b := make([]byte, n)
	nRead, err := m.r.Read(b)
	if err != nil {
		return nil, err
	}
	if nRead == n {
		return b, nil
	}
	return b[:nRead], nil
}

func (m *reader) ReadLineByDeli(deli []byte) ([]byte, error) {
	if deli == nil || len(deli) == 0 {
		return nil, fmt.Errorf("no delimiter given")
	}
	return m.r.ReadBytes(deli[0])
}

func (m *reader) Close() {
	m.f.Close()
}
