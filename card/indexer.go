package card

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/andscoop/ancli/config"
	"github.com/karrick/godirwalk"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetIndex retrieves an Index object
func GetIndex(indexPath string) (map[string]File, error) {
	index, err := loadIndex(indexPath)
	if err != nil {
		return nil, err
	}

	return index, nil
}

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
	fmt.Println(fp)
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
		fmt.Println(scanner.Text())
		if strings.Contains(scanner.Text(), "#ancli") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func saveIndex(index *map[string]File) error {
	indexData, err := json.Marshal(index)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("test.json", indexData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadIndex(fp string) (map[string]File, error) {
	indexBytes, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	var index = make(map[string]File)
	err = json.Unmarshal(indexBytes, &index)
	if err != nil {
		return nil, err
	}

	return index, nil
}

func Walk(dirname string, showHidden bool) error {
	index, err := loadIndex(dirname)
	if err != nil {
		return err
	}

	err = godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !showHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}

			// todo time cutoff logic is incorrect
			x, err := shouldIndex(osPathname)
			if err != nil {
				return err
			}

			if f, ok := index[osPathname]; ok {
				f.LastIndexed = time.Now()
				index[osPathname] = f
				return nil
			}

			if x {
				index[osPathname] = File{FilePath: osPathname, LastIndexed: time.Now()}
			}

			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})

	if err != nil {
		return err
	}

	err = saveIndex(&index)
	if err != nil {
		return err
	}

	return nil
}
