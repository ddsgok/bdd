package golden

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ddspog/str"
)

// fmtFeature takes feature name and join without spaces, adequate for
// a file name.
func fmtFeature(f string) (r string) {
	r = str.New(f).Split(" ").String()
	return
}

// filename returns the full path for the file received, with correct
// suffix at the end of file.
func filename(name string) (f string) {
	f = filepath.Join(DataDir, str.New("%s%s", name, FileNameSuffix).String())
	return
}

// getBytes returns the bytes extracted from file received.
func getBytes(name string) (bytes []byte, err error) {
	path := filename(name)
	bytes, err = ioutil.ReadFile(path)
	return
}

// ensureDir verifies if specified dir is a dir, otherwise returns an
// err, containing information about the file found.
func ensureDir(path string) (err error) {
	var info os.FileInfo
	info, err = os.Stat(path)

	switch {
	case err != nil && os.IsNotExist(err):
		err = os.MkdirAll(path, DirPerms)
	case err == nil && !info.IsDir():
		err = NewErrDataDirIsFile(path)
	}

	return
}