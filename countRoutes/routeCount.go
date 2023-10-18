package countRoutes

import (
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func RouteCount(pairs [][]string) {
	for _, pair := range pairs {
		// extract node name from the pair
		path := pair[0] // extract the first element of the pair
		// split the path by backslashes and get the last component (file name)
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		// split the file name by hyphens and get the second component (node name)
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/route table count:", nodeName)
		pterm.Info.Printf("Comparing Route Table Count: %s and %s\n", pair[1], pair[0])
		file1Content, err := os.Open(pair[1])
		if err != nil {
			// Handle error
			continue
		}
		defer file1Content.Close()

		file2Content, err := os.Open(pair[0])
		if err != nil {
			// Handle error
			continue
		}
		defer file2Content.Close()

		lines1, err := io.ReadAll(file1Content)
		if err != nil {
			// Handle error
			continue
		}

		lines2, err := io.ReadAll(file2Content)
		if err != nil {
			// Handle error
			continue
		}

		lines1Array := strings.Split(string(lines1), "\n")
		lines2Array := strings.Split(string(lines2), "\n")

		for i, line1 := range lines1Array {
			line2 := lines2Array[i]

			re := regexp.MustCompile(`^\s*Count\s*:\s*(\d+)`)

			if re.MatchString(line1) {
				// extract integer value from line1
				match1 := re.FindStringSubmatch(line1)
				if len(match1) < 2 {
					log.Fatalf("line %d in file1 doesn't contain integer after 'Count:'", i)
				}
				int1, err := strconv.Atoi(match1[1])
				if err != nil {
					log.Fatalf("line %d in file1 contains non-integer after 'Count:'", i)
				}

				// extract integer value from line2
				match2 := re.FindStringSubmatch(line2)
				if len(match2) < 2 {
					log.Fatalf("line %d in file2 doesn't contain integer after 'Count:'", i)
				}
				int2, err := strconv.Atoi(match2[1])
				if err != nil {
					log.Fatalf("line %d in file2 contains non-integer after 'Count:'", i)
				}

				if int1 == int2 {
					// integer values match
					// find the 4th line above the "Count:" should be the line with the command
					var lineBeforeCount string
					if i >= 4 {
						lineBeforeCount = lines1Array[i-4]

					} else {
						lineBeforeCount = ""
					}
					// split the line to extract only the command
					lineSplit := strings.Split(string(lineBeforeCount), ";")
					pterm.Success.Println("Route count matches ->", lineSplit[2])
				} else {
					// find the 4th line above the "Count:" should be the line with the command
					var lineBeforeCount string
					if i >= 4 {
						lineBeforeCount = lines1Array[i-4]

					} else {
						lineBeforeCount = ""
					}
					// split the line to extract only the command
					lineSplit := strings.Split(string(lineBeforeCount), ";")
					// integer values don't match
					pterm.Warning.Println("Route count doesn't match ->", lineSplit[2])
					// print non-matching lines
					pterm.FgRed.Printf("Line %d in %s %s\n", i, pair[1], line1)
					pterm.FgRed.Printf("Line %d in %s %s\n", i, pair[0], line2)
				}
			}
		}
	}
}
