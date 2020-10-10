package deck

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestCardIndexer(t *testing.T) {
	var tests = []struct {
		name            string
		relFP           string
		cardRegex       string
		fpRegex         string
		wantShouldIndex bool
	}{
		{"DontIndexDirs", "/examples/", "*", ".md", false},
		{"IndexExactMatch", "/examples/band.md", "#ancli-jokes", ".md", true},
		{"IndexRegexMatch1", "/examples/band.md", "#ancli*", ".md", true},
		{"WrongFileType", "/examples/band.md", "*", ".x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			path := strings.SplitN(wd, "/", -1)
			rootDir := strings.Join(path[0:len(path)-1], "/")

			haveShouldIndex, err := shouldIndex(rootDir+tt.relFP, tt.cardRegex, tt.fpRegex)
			if err != nil {
				panic(err)
			}

			if haveShouldIndex != tt.wantShouldIndex {
				t.Errorf("expected shouldIndex=%v but got shouldIndex=%v", tt.wantShouldIndex, haveShouldIndex)
			}
		})
	}
}
