package deck

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
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

func shouldIndex(fp, deckRegex, fpRegex string) (bool, error) {
	if isDir(fp) {
		return false, nil
	}

	fpMatch, err := regexp.MatchString(fpRegex, fp)
	if !fpMatch {
		return false, nil
	}

	file, err := os.Open(fp)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		match, _ := regexp.MatchString(deckRegex, scanner.Text())
		if match {
			return true, nil
		}
	}

	return false, nil
}

// IndexAndSave rebuilds the deck index and saves it to a file
func (d *Deck) IndexAndSave(indexHidden bool) error {
	err := d.Index(indexHidden)
	if err != nil {
		return err
	}

	d.Save()

	return nil
}

// Index walks a directory looking for cards to add to a deck
func (d *Deck) Index(indexHidden bool) error {
	cards := make(map[string]*Card)
	err := godirwalk.Walk(d.RootDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			osPathname = strings.ToLower(osPathname)
			fpHash := hashFp(osPathname)

			if !indexHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}

			fpRegex := config.GetString("fpRegex")

			x, err := shouldIndex(osPathname, d.DeckRegex, fpRegex)
			if err != nil {
				return err
			}

			// Update existing cards
			if c, ok := cards[fpHash]; ok {
				c.LastIndexed = time.Now().Format(time.RFC3339)
				c.IsArchived = false
				cards[fpHash] = c
			} else {
				if x {
					cards[fpHash] = &Card{Fp: osPathname, LastIndexed: time.Now().Format(time.RFC3339), IsArchived: false}
				}
			}

			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})

	d.Cards = cards
	if err != nil {
		return err
	}

	d.LastIndexed = time.Now().Format(time.RFC3339)
	return nil
}
