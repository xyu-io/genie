package queue

import (
	"encoding/json"
	"fmt"
	"github.com/xyu-io/genie/looper"
	"os"
	"path/filepath"
	"time"
)

// DumperReadFn is function to read dump file
type DumperReadFn func(mQueue Queue, data []byte) error

// Dumper is queue support dump
type Dumper interface {
	SetDir(dir string)
	Read(interval int, fn DumperReadFn)
	ReadOnce(fn DumperReadFn) error
	Dump(size int) (int, error)
	Queue
}

type baseDumper struct {
	Queue
	// 导出的目录，里面可能包含多个文件
	dir string
	// 队列大小
	size int
}

// SetDir sets the dump dir
func (bd *baseDumper) SetDir(dir string) {
	bd.dir = dir
	_ = os.MkdirAll(bd.dir, os.ModePerm)
}

// NewDumper returns default dumper queue
func NewDumper(size int, dir string) Dumper {
	bd := &baseDumper{
		Queue: New(size),
		size:  size,
	}
	bd.SetDir(dir)
	return bd
}

// filename returns the dump filename
func (bd *baseDumper) filename() string {
	return fmt.Sprintf("data_%s.dump", time.Now().Format("060102150405"))
}

func (bd *baseDumper) Dump(size int) (int, error) {
	if size == 0 {
		size = bd.size * 2
	}
	data := bd.PopBatch(size)
	if len(data) == 0 {
		return 0, nil
	}
	fileBytes, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	file := filepath.Join(bd.dir, bd.filename())
	return len(data), os.WriteFile(file, fileBytes, 0644)
}

// Read reads the data from dump dir cyclically
func (bd *baseDumper) Read(interval int, fn DumperReadFn) {
	looper.TimeLoopThen(time.Second*time.Duration(interval), false, func(now time.Time) {
		_ = bd.ReadOnce(fn)
	})
}

// ReadOnce reads the data from dump dir recursively
func (bd *baseDumper) ReadOnce(fn DumperReadFn) error {
	var (
		isRead bool
	)
	return filepath.Walk(bd.dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 一次只读一个文件
		if isRead || info.IsDir() || filepath.Ext(fpath) != ".dump" {
			return nil
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			if err := fn(bd.Queue, data); err != nil {
				return err
			}
		}
		_ = os.RemoveAll(fpath)
		isRead = true
		return nil
	})
}
