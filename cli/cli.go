package cli

import (
	"fmt"
	hubbub "github.com/rjz/hubbub/common"
	"os"
	"path/filepath"
)

func die(message string, err error) {
	fmt.Println("failed reading policies directory")
	fmt.Println(err)
	os.Exit(1)
}

func prettyListItem(listItem string) {
	fmt.Println("  *", listItem)
}

func prettyListHeader(header string) {
	fmt.Println("\n", header)
}

func prettyListFooter() {
	fmt.Println("")
}

func prettyTable(title string, items []string) {
	prettyListHeader(title)
	for _, v := range items {
		prettyListItem(v)
	}
	prettyListFooter()
}

func ListConfigFiles(path string) {
	ext := ".json"
	pathGlob := filepath.Join(path, fmt.Sprintf("*%s", ext))
	matches, err := filepath.Glob(pathGlob)
	if err != nil {
		die("failed reading config directory", err)
	}
	var items []string
	for _, v := range matches {
		basename := filepath.Base(v)
		items = append(items, basename[0:len(basename)-len(ext)])
	}
	prettyTable(path, items)
}

// printGoals describes all globally-registered goals
func ListPolicyGoals() {
	serviceFactories := hubbub.ServiceFactories()
	prettyTable("goals", serviceFactories.Goals())
}
