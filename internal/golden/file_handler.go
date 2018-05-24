package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrFileNotFoundInDataDir is thrown when the file to open isn't
	// found on data dir.
	ErrFileNotFoundInDataDir = errors.New("file not founded in data dir")
)

// fileHandler stores full file path, and can perform some operations.
type fileHandler string

// Path returns the full path for file.
func (fh fileHandler) Path() (s string) {
	s = string(fh)
	return
}

// Dir returns the full path for current dir of file.
func (fh fileHandler) Dir() (d string) {
	d = filepath.Dir(fh.Path())
	return
}

// EnsureDir verifies if file dir is a directory, otherwise returns an
// err, containing information about the file found. If it doesn't
// exist, it will create the directory.
func (fh fileHandler) EnsureDir() (err error) {
	var info os.FileInfo
	info, err = os.Stat(fh.Dir())

	switch {
	case err != nil && os.IsNotExist(err):
		err = os.MkdirAll(fh.Dir(), DirPerms)
	case err == nil && !info.IsDir():
		err = NewErrDataDirIsFile(fh.Dir())
	}

	return
}

// Bytes return the file content as bytes.
func (fh fileHandler) Bytes() (bytes []byte, err error) {
	bytes, err = ioutil.ReadFile(fh.Path())
	return
}

// Write put bytes in file content.
func (fh fileHandler) Write(bytes []byte) (err error) {
	err = ioutil.WriteFile(fh.Path(), bytes, FilePerms)
	return
}

// Ext returns the extension of the file.
func (fh fileHandler) Ext() (e string) {
	e = filepath.Ext(fh.Path())
	return
}

// ExtType returns the type of file.
func (fh fileHandler) ExtType() (et fileType) {
	et = fileType(fh.Ext())
	return
}

// path look for a file with name received inside data dir, when found
// it will return a file handler with information about the file. If
// more than one file exists with same, it will use the first found.
func path(name string) (fh fileHandler, err error) {
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(DataDir); err == nil {
		for _, file := range files {
			if strings.Contains(file.Name(), name) {
				fh = fileHandler(
					filepath.Join(
						DataDir,
						file.Name(),
					),
				)

				return
			}
		}
	}

	err = ErrFileNotFoundInDataDir

	return
}
