package uClig

import (
	"fmt"
	"strings"
)

var (
	Debug = true
)

// DisplayList is a helper function that for printing the lists generated in this file.
// As a simple implementation, string splitting by ", " economizes the []string type
// into a [][]string type without introducing more panics or complicated type assertions.
func DisplayList(sn string, list []string) {
	fmt.Printf("SN: %s", sn)
	for i := range list {
		splitEntry := strings.Split(list[i], ", ")
		for _, e := range splitEntry {
			fmt.Printf("\n\t%s", e)
		}
	}
	fmt.Println()
}
