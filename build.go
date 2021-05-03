/*
  Tool for compiling small snippets in `nsroot` directory into a single snippet file
	at `snippets/snippets` and generating `COMMANDS.md`.

	To compile  to binary (you need golang installed):
  $ go build build.go

  To run without build (you need golang installed):
  $ go run build.go
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

const Version = "1.1.0" // 2021-05-02
const root = "nsroot"
const tabSpace = "  "

type Snippet struct {
	Prefix      interface{} `json:"prefix"`
	Body        interface{} `json:"body"`
	Description string      `json:"description"`
}

func main() {

	if len(os.Args) > 1 { // No argument accepted
		fmt.Printf(`Shellman build tool v%v

This tool doesn't accept any argument. Run it from project root directory.
It concatenates 'nsroot' snippets to 'snippets/snippets.json' and generates 'COMMANDS.md'.
`, Version)
		os.Exit(1)
	}

	cwd, errGetCurrentDirectory := os.Getwd()
	if errGetCurrentDirectory != nil {
		fmt.Printf("Cannot get current directory due to error: %v\n", errGetCurrentDirectory)
		panic(errGetCurrentDirectory)
	}
	snippetFilePath := filepath.Join(cwd, "snippets", "snippets.json")
	snippetFilePath = filepath.FromSlash(snippetFilePath) // OS agnostic path
	docFilePath := filepath.Join(cwd, "COMMANDS.md")
	docFilePath = filepath.FromSlash(docFilePath) // OS agnostic path

	docTitleBuilder := strings.Builder{}
	docBodyBuilder := strings.Builder{}

	docTitleBuilder.WriteString("# Commands\n\n")

	changeToRootErr := os.Chdir(root)
	if changeToRootErr != nil {
		fmt.Println("Run this program from Shellman root directory")
		panic(changeToRootErr)
	}

	folders := getFolders(".")
	snippets := map[string]Snippet{}

	for _, folder := range folders {
		docTitleBuilder.WriteString("### " + folder + "\n\n")
		files := getFiles(folder)

		errChangeToChildDir := os.Chdir(folder)
		if errChangeToChildDir != nil {
			fmt.Printf("Cannot change path to child directory: %s due to error: %v\n",
				folder, errChangeToChildDir)
			panic(errChangeToChildDir)
		}

		for _, file := range files {

			snippetName := folder + "." + strings.TrimSuffix(file, ".json")
			jsonFileBytes, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Cannot read file: %s due to error: %v\n", file, err)
				panic(err)
			}

			var snippet Snippet
			json.Unmarshal(jsonFileBytes, &snippet)
			namedSnippet := map[string]Snippet{
				snippetName: snippet,
			}

			for k, v := range namedSnippet {
				snippets[k] = v

				prefixes := []string{}
				docTitleBuilder.WriteString("  - [")
				docBodyBuilder.WriteString("## ")

				switch v.Prefix.(type) {
				case []interface{}: // snippet with multiple prefixes
					for _, prefix := range v.Prefix.([]interface{}) {
						p := prefix.(string)
						prefixes = append(prefixes, p)
					}
					link := strings.Join(prefixes[:], " , ")
					docTitleBuilder.WriteString(link)
					docTitleBuilder.WriteString("](#")

					docBodyBuilder.WriteString(link + "\n\n")
					docBodyBuilder.WriteString(v.Description + "[&uarr;](#" + folder + ")\n\n")

					link = strings.ReplaceAll(link, " ", "-")
					docTitleBuilder.WriteString(link)
				case interface{}: // snippet with single prefix
					p := v.Prefix.(string)
					docTitleBuilder.WriteString(p)
					docTitleBuilder.WriteString("](#")

					docBodyBuilder.WriteString(p + "\n\n")
					docBodyBuilder.WriteString(v.Description + "[&uarr;](#" + folder + ")\n\n")

					link := strings.ReplaceAll(p, " ", "-")
					docTitleBuilder.WriteString(link)
				default:
					panic("Unknown prefix type")
				}

				docTitleBuilder.WriteString(")\n\n")
			}
		}
		errChangeToParentDir := os.Chdir("..")
		if errChangeToParentDir != nil {
			fmt.Printf("Cannot change path to parent directory: %s due to error: %v\n",
				folder, errChangeToParentDir)
			panic(errChangeToParentDir)
		}
	}

	// write snippet file
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", tabSpace)
	jsonEncoder.Encode(snippets)
	snippetFileErr := ioutil.WriteFile(snippetFilePath, buf.Bytes(), 0755)
	if snippetFileErr != nil {
		fmt.Printf("Cannot write snippet file: %v due to error: %v\n",
			snippetFilePath, snippetFileErr)
		panic(snippetFileErr)
	}

	// Write documentation file
	docFile, docFileErr := os.Create(docFilePath)
	if docFileErr != nil {
		fmt.Printf("Cannot create documentation file: %v due to error: %v\n",
			docFilePath, docFileErr)
		panic(docFileErr)
	}
	defer docFile.Close()
	documentation := docTitleBuilder.String() + docBodyBuilder.String()
	_, docFileWriteErr := docFile.WriteString(documentation)
	if docFileWriteErr != nil {
		fmt.Printf("Cannot write documentation file: %v due to error: %v\n",
			docFilePath, docFileWriteErr)
		panic(docFileWriteErr)
	}

}

func getFolders(root string) []string {
	var folders []string
	files, err := ioutil.ReadDir(root)
	if err != nil {
		fmt.Printf("Cannot read "+root+" directories: %v due to error: %v\n", root, err)
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
		fmt.Printf("Cannot read directory: %v due to error: %v\n", dir, err)
		panic(err)
	}

	var result []string
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			fmt.Printf("Namespace directories should not contain nested directories: '%v' contains '%v'\n", dir, fileInfo.Name())
			panic("Found nested directories inside namespace directory")
		}
		result = append(result, fileInfo.Name())
	}
	return result
}
