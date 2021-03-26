/*
  Tool for compiling small snippets in `nsroot` directory into a single snippet file
	at `snippets/snippets`. To compile (you need golang installed):
  # go build build.go
  To run without build
  # go run build.go

	Date:    2021-03-26
	Version: 1.0.0
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const root = "nsroot"
const tabSpace = "  "

type Snippet struct {
	Prefix      interface{} `json:"prefix"`
	Body        interface{} `json:"body"`
	Description string      `json:"description"`
}

func main() {
	cwd, errGetCurrentDirectory := os.Getwd()
	if errGetCurrentDirectory != nil {
		fmt.Printf("Cannot get current directorydue to error: %v", errGetCurrentDirectory)
		panic(errGetCurrentDirectory)
	}
	snippetFile := filepath.Join(cwd, "snippets", "snippets.json")
	snippetFile = filepath.FromSlash(snippetFile) // OS agnostic path

	os.Chdir(root)
	folders := getFolders(".")
	snippets := map[string]Snippet{}

	for _, folder := range folders {
		files := getFiles(folder)
		for _, file := range files {
			errChangeToChildDir := os.Chdir(folder)
			if errChangeToChildDir != nil {
				fmt.Printf("Cannot change path to child directory: %v due to error: %v", folder, errChangeToChildDir)
				panic(errChangeToChildDir)
			}
			snippetName := folder + "." + strings.TrimSuffix(file, ".json")
			jsonFileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Cannot read file: %v due to error: %v", file, err)
				panic(err)
			}

			var snippet Snippet
			json.Unmarshal(jsonFileBytes, &snippet)
			namedSnippet := map[string]Snippet{
				snippetName: snippet,
			}

			for k, v := range namedSnippet {
				snippets[k] = v
			}
			errChangeToParentDir := os.Chdir("..")
			if errChangeToParentDir != nil {
				fmt.Printf("Cannot change path to parent directory: %v due to error: %v", folder, errChangeToParentDir)
				panic(errChangeToParentDir)
			}
		}
	}
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", tabSpace)
	jsonEncoder.Encode(snippets)
	err := ioutil.WriteFile(snippetFile, buf.Bytes(), 0755)
	if err != nil {
		fmt.Printf("Cannot write snippet file: %v due to error: %v", snippetFile, err)
		panic(err)
	}
}

func getFolders(root string) []string {
	var folders []string
	files, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Printf("Cannot read "+root+" directories: %v due to error: %v", root, err)
		panic(err)
	}
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			folders = append(folders, fileInfo.Name())
		}
	}
	return folders
}

func getFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Cannot read directory: %v due to error: %v", dir, err)
		panic(err)
	}

	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result
}
