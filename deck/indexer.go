package deck

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

func shouldIndex(fp, deckPrefix string) (bool, error) {
	ext := config.GetString("cardFileExt")
	fpParts := strings.Split(fp, ".")

	if fpParts[len(fpParts)-1] != ext {
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
		if strings.Contains(scanner.Text(), deckPrefix) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// Walk takes a stroll through the dir
// skips hidden files by default
func Walk(dir string, showHidden bool) error {
	d := NewDeck()

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			osPathname = strings.ToLower(osPathname)

			fmt.Println(osPathname)
			if !showHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}

			// todo time cutoff logic is incorrect
			deckPrefix := config.GetString("deckPrefix")

			x, err := shouldIndex(osPathname, deckPrefix)
			if err != nil {
				return err
			}

			// Update existing card
			if c, ok := d.Cards[osPathname]; ok {
				c.LastIndexed = time.Now().String()
				d.Cards[osPathname] = c
			} else {
				if x {
					d.Cards[osPathname] = Card{Fp: osPathname, LastIndexed: time.Now().String()}
				}
			}

			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})

	err = d.UpdateKeys()
	if err != nil {
		return err
	}

	fmt.Print(d.Cards)

	config.SetAndSave("decks", d.Cards)
	return nil
}
