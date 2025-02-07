package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// CHange node ID in communities.txt to the corresponding title in amazon-meta.txt

	f, err := os.Open("amazon-meta.txt")
	if err != nil {
		fmt.Println("Error opening file")
		return
	}
	defer f.Close()
	IdTitleMap := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Id:") {
			id := line[6:]
			// Skip the next 2 lines
			scanner.Scan()
			scanner.Scan()
			line = scanner.Text()
			if strings.HasPrefix(line, "  title: ") {
				title := line[9:]
				IdTitleMap[id] = title
			}
		}
	}

	f2, err := os.Open("communities.txt")
	if err != nil {
		fmt.Println("Error opening file")
		return
	}
	f3, err := os.Create("communitiesTitle.txt")
	if err != nil {
		fmt.Println("Error creating file")
		return
	}

	defer f2.Close()
	scanner2 := bufio.NewScanner(f2)
	for scanner2.Scan() {
		line := scanner2.Text()
		if strings.HasPrefix(line, "Node") {
			words := strings.Fields(line)
			id := words[1]
			title := IdTitleMap[id]
			f3.WriteString("	" + title + "\n")
		}
		if strings.HasPrefix(line, "Community") {
			f3.WriteString(line + "\n")
		}

	}
}
