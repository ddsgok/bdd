package golden

import (
	"fmt"
)

// ErrDataDirIsFile is an error thrown when the data folder, with the
// golden files is a file.
type ErrDataDirIsFile struct {
	file string
}

// Error prints the reason for error, telling which file was found.
func (e *ErrDataDirIsFile) Error() (m string) {
	m = fmt.Sprintf("data folder is a file: %s", e.file)
	return
}

// NewErrDataDirIsFile creates an ErrDataDirIsFile with the file founded.
func NewErrDataDirIsFile(file string) (err *ErrDataDirIsFile) {
	err = &ErrDataDirIsFile{file: file}
	return
}