package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cjun714/glog/log"
)

func main() {
	path := os.Args[1]

	dirs, _ := os.ReadDir(path)

	for _, fi := range dirs {
		if !fi.IsDir() {
			continue // skip non-dir
		}

		log.I("packing:", fi.Name())
		if e := zipDir(filepath.Join(path, fi.Name())); e != nil {
			log.E("pack failed:", fi.Name(), e)
		}
	}
}

func zipDir(src string) error {
	target := src + ".cbz"
	f, e := os.Create(target)
	if e != nil {
		return e
	}
	defer f.Close()

	wr := zip.NewWriter(f)
	defer wr.Close()

	e = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !isImage(info.Name()) {
			log.E("non-image:", info.Name())
			return nil
		}

		data, e := ioutil.ReadFile(path)
		if e != nil {
			return e
		}

		w, e := wr.Create(info.Name())
		if e != nil {
			return e
		}
		_, e = w.Write(data)

		return e
	})

	return e
}

func isImage(name string) bool {
	ext := filepath.Ext(name)
	ext = strings.ToLower(ext)

	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" ||
		ext == ".webp" || ext == ".gif" || ext == ".bmp"
}
