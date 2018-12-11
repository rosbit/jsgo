package mod_fs

/**
 * fs module implementation, which is invoked via `var fs = require('fs')` in js.
 * NOTE: This module is different from Node.js fs module.
 * Rosbit Xu <me@rosbit.cn>
 * Dec. 9, 2018
 */

import (
	js "github.com/rosbit/duktape-bridge/duk-bridge-go"
	"os"
	"io/ioutil"
	"fmt"
)

type FsModule struct {
	fds map[int]*os.File
}

func NewFsModule(ctx *js.JSEnv) interface{} {
	return &FsModule{make(map[int]*os.File)}
}

func stringToPosixFlags(flags string) int {
	posixFlags := 0
	switch flags {
	case "r":
		posixFlags = os.O_RDONLY
	case "r+":
		posixFlags = os.O_RDWR
	case "w":
		posixFlags = os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	case "w+":
		posixFlags = os.O_RDWR | os.O_TRUNC | os.O_CREATE
	case "a":
		posixFlags = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	case "a+":
		posixFlags = os.O_RDWR | os.O_APPEND | os.O_CREATE
	}
	return posixFlags
}

func (m *FsModule) Open(path string, flags string, mode uint32) (int, error) {
	var f *os.File
	var e error
	if flags == "" {
		f, e = os.Open(path)
	} else {
		posixFlags := stringToPosixFlags(flags)
		var perm os.FileMode
		if mode == 0 {
			perm = 0666
		} else {
			perm = os.FileMode(mode) & os.ModePerm
		}
		f, e = os.OpenFile(path, posixFlags, perm)
	}

	if e != nil {
		return -1, e
	}
	fd := int(uint64(f.Fd()))
	m.fds[fd] = f
	return fd, nil
}

func (m *FsModule) Close(fd int) {
	if f, ok := m.fds[fd]; ok {
		f.Close()
		delete(m.fds, fd)
	}
}

func (m *FsModule) Write(fd int, data []byte, position int64) (int, error) {
	f, ok := m.fds[fd]
	if !ok {
		return -1, fmt.Errorf("fd %d is invalid", fd)
	}

	var n int
	var e error
	if position <= 0 {
		n, e = f.Write(data)
	} else {
		n, e = f.WriteAt(data, position)
	}
	if e != nil {
		return -2, e
	}
	return n, nil
}

func (m *FsModule) Read(fd int, length int, position int64) ([]byte, error) {
	f, ok := m.fds[fd]
	if !ok {
		return nil, fmt.Errorf("fd %d is invalid", fd)
	}
	if length <= 0 {
		return nil, nil
	}

	data := make([]byte, length)
	var n int
	var e error
	if position <= 0 {
		n, e = f.Read(data)
	} else {
		n, e = f.ReadAt(data, position)
	}
	if e != nil {
		return nil, e
	}
	if n == 0 {
		// eof
		return nil, nil
	}
	if n == length {
		return data, nil
	}
	return data[:n], nil
}

func (m *FsModule) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (m *FsModule) WriteFile(filename string, data []byte) (int, error) {
	if filename == "" {
		return -1, fmt.Errorf("filename expected")
	}
	if data == nil || len(data) == 0 {
		return 0, nil
	}
	err := ioutil.WriteFile(filename, data, 0755)
	if err != nil {
		return -2, err
	}
	return len(data), nil
}

func (m *FsModule) AppendFile(filename string, data []byte) (int, error) {
	if filename == "" {
		return -1, fmt.Errorf("filename expected")
	}
	if data == nil {
		return 0, nil
	}
	f, e := os.OpenFile(filename, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0755)
	if e != nil {
		return -2, e
	}
	defer f.Close()

	n, e := f.Write(data)
	if e != nil {
		return -3, e
	}
	return n, e
}
