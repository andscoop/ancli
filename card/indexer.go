package card

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andscoop/ancli/config"
	"github.com/karrick/godirwalk"
)

func isHidden(osPathname string) bool {
	fileNodes := strings.Split(osPathname, "/")
	node := fileNodes[len(fileNodes)-1]

	if string(node[0]) == "." {
		return true
	}

	return false
}

func isDir(fp string) bool {
	info, err := os.Stat(fp)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func shouldIndex(fp string) (bool, error) {
	// todo temp skip scanning the go binary
	if fp == "/Users/andrew.cooper/go/src/github.com/andscoop/ancli/ancli" {
		return false, nil
	}

	if isDir(fp) {
		return false, nil
	}

	file, err := os.Open(fp)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "#ancli") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func Walk(dirname string, showHidden bool) error {
	index, err := config.GetIndex()
	if err != nil {
		return err
	}

	fmt.Println(index)
	err = godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			osPathname = strings.ToLower(osPathname)

			if !showHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}

			// todo time cutoff logic is incorrect
			x, err := shouldIndex(osPathname)
			if err != nil {
				return err
			}

			if f, ok := index[osPathname]; ok {
				f.LastIndexed = time.Now().String()
				index[osPathname] = f
				return nil
			}

			if x {
				index[osPathname] = config.Index{FilePath: osPathname, LastIndexed: time.Now().String()}
			}

			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})

	if err != nil {
		return err
	}

	//err = saveIndex(&index)
	//if err != nil {
	//	return err
	//}
	c := config.GetConfig()
	c.Set("decks", index)
	config.SaveConfig(c)
	return nil
}
