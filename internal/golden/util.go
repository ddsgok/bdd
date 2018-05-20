package golden

import (
	"github.com/ddspog/str"
	"io/ioutil"
	"os"
	"path/filepath"
)

func fmtFeature(f string) (r string) {
	r = str.New(f).Split(" ").String()
	return
}
func filename(name string) (f string) {
	f = filepath.Join(DataDir, str.New("%s%s", name, FileNameSuffix).String())
	return
}
func getBytes(name string) (bytes []byte, err error) {
	path := filename(name)
	bytes, err = ioutil.ReadFile(path)
	return
}
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
