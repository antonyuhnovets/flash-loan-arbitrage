package Filestorage

import (
	c "context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type FileStorage struct {
	Files map[string]string
}

func NewStorage() (
	fs *FileStorage,
) {
	fs = &FileStorage{
		Files: make(map[string]string),
	}

	return
}

func (fs *FileStorage) Setup(
	Files map[string]string,
) (
	err error,
) {
	for k, v := range Files {
		err = fs.NewFile(k, v)
		if err != nil {
			return
		}
	}

	return
}

func (fs *FileStorage) UseFile(
	name, path string,
) {
	fs.Files[name] = path
}

func (fs *FileStorage) NewFile(
	name, path string,
) (

	err error,
) {
	f, err := os.Create(path)

	if err != nil {
		return
	}
	defer f.Close()

	fs.UseFile(
		name,
		path,
	)

	fs.Store(
		c.Background(),
		name,
		[]byte("[\n"),
	)
	err = fs.Store(
		c.Background(),
		name,
		[]byte("]"),
	)

	return
}

func (fs *FileStorage) Store(
	ctx c.Context,
	where string,
	item interface{},
) (
	err error,
) {
	f, err := os.OpenFile(
		fs.Files[where],
		os.O_APPEND|os.O_RDWR,
		0644,
	)
	if err != nil {
		return
	}

	n, err := f.Write(item.([]byte))
	log.Println(n)

	return
}

func (fs *FileStorage) Read(
	ctx c.Context,
	where string,
	out interface{},
) (
	err error,
) {
	b, err := os.ReadFile(
		fs.Files[where],
	)
	if err != nil {
		return
	}
	_, ok := out.(*[]byte)
	if !ok {
		err = json.Unmarshal(b, out)
		if err != nil {

			return
		}
	} else {
		res := &b
		out = res
	}

	return
}

func (fs *FileStorage) ContinueFile(
	ctx c.Context,
	where string,
	item interface{},
) (
	err error,
) {
	f, err := os.OpenFile(
		fs.Files[where],
		os.O_RDWR,
		0644,
	)
	if err != nil {
		return
	}

	b, err := os.ReadFile(fs.Files[where])
	if err != nil {
		log.Println(err)

		return
	}

	if len(b) < 3 {
		b = []byte("[\n]")
	}
	if string(b[len(b)-3]) != "[" {
		b = append(b[:len(b)-2], []byte(",\n")...)
	} else {
		b = append(b[:len(b)-2], []byte("\n")...)
	}

	b = append(b, item.([]byte)...)

	n, err := f.Write(b)
	log.Println(n)

	return
}

func (fs *FileStorage) Remove(
	ctx c.Context,
	where string,
	item interface{},
) (
	err error,
) {
	f, err := os.OpenFile(
		fs.Files[where],
		os.O_RDWR,
		0644,
	)
	if err != nil {
		return
	}
	defer f.Close()

	b := make([]byte, 0)

	err = fs.Read(ctx, where, &b)
	if err != nil {

		return
	}

	before, after, ok := strings.Cut(string(b), string(item.([]byte)))
	if !ok {
		err = fmt.Errorf("item %b not found", item)

		return
	}

	out := []byte(before[:len(before)-2] + after)

	fs.Clear(ctx, where)

	n, err := f.Write(out)

	fmt.Println("written bytes ", n)

	return
}

func (fs *FileStorage) Clear(
	ctx c.Context,
	where string,
) (
	err error,
) {
	err = os.Truncate(
		fs.Files[where],
		0,
	)

	return
}

func (fs *FileStorage) ClearAll(
	ctx c.Context,
) (
	err error,
) {
	for k, v := range fs.Files {
		err = os.Remove(v)
		if err != nil {
			return
		}
		fs.Files[k] = ""
	}

	return
}
