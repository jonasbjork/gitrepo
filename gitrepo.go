package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	gitURLs := findGitRepos(baseDir)
	if len(gitURLs) > 0 {
		for _, url := range gitURLs {
			fmt.Println(url)
		}
	} else {
		fmt.Println("No git repositories found in", baseDir)
	}
}

func findGitRepos(baseDir string) []string {
	var gitURLs []string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			configPath := filepath.Join(path, "config")
			gitURL := getGitURL(configPath)
			if gitURL != "" {
				gitURLs = append(gitURLs, gitURL)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}

	return gitURLs
}

func getGitURL(configPath string) string {
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isRemoteOrigin := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "[remote \"origin\"]") {
			isRemoteOrigin = true
		}

		if isRemoteOrigin && strings.HasPrefix(line, "url =") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading config file:", err)
	}

	return ""
}
