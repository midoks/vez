// Copyright 2020 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/flamego/template"
	"github.com/midoks/vez/internal/tools"
)

//go:generate go-bindata -nomemcopy -nometadata -ignore="\\.DS_Store" -pkg=templates -prefix=../../../templates -debug=false -o=templates_gen.go ../../../templates/...

type fileSystem struct {
	files []template.File
}

func (fs *fileSystem) Files() []template.File { return fs.files }

func (fs *fileSystem) Get(name string) (io.Reader, error) {
	for i := range fs.files {
		if fs.files[i].Name()+fs.files[i].Ext() == name {
			bReader, _ := fs.files[i].Data()
			return bytes.NewReader(bReader), nil
		}
	}
	return nil, fmt.Errorf("file %q not found", name)
}

type file struct {
	name string
	data []byte
	ext  string
}

func (f *file) Name() string          { return f.name }
func (f *file) Data() ([]byte, error) { return f.data, nil }
func (f *file) Ext() string           { return f.ext }

// NewTemplateFileSystem returns a template.TemplateFileSystem instance for embedded assets.
// The argument "dir" can be used to serve subset of embedded assets. Template file
// found under the "customDir" on disk has higher precedence over embedded assets.
func NewTemplateFileSystem(dir, customDir string) template.FileSystem {

	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	var files []template.File
	names := AssetNames()
	for _, name := range names {
		if !strings.HasPrefix(name, dir) {
			continue
		}

		// Check if corresponding custom file exists
		var err error
		var data []byte
		fpath := path.Join(customDir, name)
		if tools.IsFile(fpath) {
			data, err = ioutil.ReadFile(fpath)
		} else {
			data, err = Asset(name)
		}
		if err != nil {
			continue
		}

		name = strings.TrimPrefix(name, dir)
		ext := path.Ext(name)
		name = strings.TrimSuffix(name, ext)

		// fmt.Println(name, data, ext)
		files = append(files, &file{
			name: name,
			data: data,
			ext:  ext,
		})
	}
	return &fileSystem{files: files}
}
