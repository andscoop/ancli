package indexer

import (
	"fmt"
	"github.com/karrick/godirwalk"
	"path/filepath"
	"strings"
	"time"
)

type card struct {
	fp           string
	lastIndexed  time.Time
	firstIndexed time.Time
}

type index struct {
	lastIndexed time.Time
	cards       []card
}

func isHidden(osPathname string) bool {
	fileNodes := strings.Split(osPathname, "/")
	node := fileNodes[len(fileNodes)-1]

	if string(node[0]) == "." {
		fmt.Println("true")
		return true
	}

	return false
}

func Walk(dirname string, showHidden bool) {
	i := index{lastIndexed: time.Now(), cards: make([]card, 1)}
	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if !showHidden && isHidden(osPathname) {
				return filepath.SkipDir
			}

			// time cutoff logic is incorrect
			c := card{fp: osPathname, lastIndexed: time.Now()}
			i.cards = append(i.cards, c)

			fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", i)
}
