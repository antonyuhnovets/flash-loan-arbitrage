package Filestorage

import (
	c "context"
	"log"
	"os"
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
	item []byte,
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

	n, err := f.Write(item)
	log.Println(n)

	return
}

func (fs *FileStorage) Read(
	ctx c.Context,
	where string,
) (
	b []byte,
	err error,
) {
	b, err = os.ReadFile(
		fs.Files[where],
	)

	return
}

func (fs *FileStorage) ContinueFile(
	ctx c.Context,
	where string,
	item []byte,
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
	b, err := fs.Read(ctx, where)
	if err != nil {
		return
	}

	if string(b[len(b)-3]) != "[" {
		b = append(b[:len(b)-2], []byte(",\n")...)
	} else {
		b = append(b[:len(b)-2], []byte("\n")...)
	}
	b = append(b, item...)

	n, err := f.Write(b)
	log.Println(n)

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
