package golden

import (
	"github.com/ddspog/str"
)

type ErrDataDirIsFile struct {
	file string
}

func NewErrDataDirIsFile(file string) (err *ErrDataDirIsFile) {
	err = &ErrDataDirIsFile{file: file}
	return
}
func (e *ErrDataDirIsFile) Error() (m string) {
	m = str.New("data folder is a file: %s", e.file).String()
	return
}
