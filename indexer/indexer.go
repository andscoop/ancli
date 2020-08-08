package indexer

import (
	"fmt"
	"github.com/karrick/godirwalk"
	"path/filepath"
)

func isHidden(osPathname string) bool {
	if string(osPathname[0]) == "." {
		return true
	}

	return false
}

func Walk(dirname string, showHidden bool) {
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !showHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}
			fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		panic(err)
	}
}
