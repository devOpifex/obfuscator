package main

import (
	"bufio"
	"log"
	"os"
)

func readLicense(path string) string {
	if path == "" {
		return ""
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open -header file:", err)
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		content += "# " + line + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading -header file:", err)
	}

	// Remove the last newline if content is not empty
	if len(content) > 0 {
		content = content[:len(content)-1]
	}

	return content + "\n"
}
