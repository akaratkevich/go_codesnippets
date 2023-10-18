package searchDirectory

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func DirSearch(patterns []string) [][]string {
	folder := "D:\\vdcMigrations"

	files, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	var pairs [][]string
	for _, pattern := range patterns {
		tempPairs := make(map[string][]string)
		for _, file := range files {
			filename := file.Name()

			// Use the regular expression pattern to match the filename
			if matched, _ := regexp.MatchString(pattern, filename); matched {
				name := filename[:5]

				// Add the file to the pair based on its name
				if _, ok := tempPairs[name]; !ok {
					tempPairs[name] = []string{filepath.Join(folder, filename)}
				} else {
					tempPairs[name] = append(tempPairs[name], filepath.Join(folder, filename))
				}
			}
		}

		// Add the pairs for the current pattern to the main pairs slice
		for _, pair := range tempPairs {
			if len(pair) >= 2 { // TODO: needs to test for logic with multiple files
				pairs = append(pairs, pair)
			}
		}
	}

	return pairs
}
