package golden

import (
	"flag"
	"os"
)

var (
	DataDir                    = "testdata"
	FileNameSuffix             = ".json"
	FlagName                   = "update"
	FilePerms      os.FileMode = 0644
	DirPerms       os.FileMode = 0755
	update                     = flag.Bool(FlagName, false, "Update golden test file fixture")
)

func WillUpdate() (b bool) {
	b = *update
	return
}
