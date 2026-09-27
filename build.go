/*
Tool for compiling the JSON files below nsroot into snippets/snippets.json and
generating COMMANDS.md.

To compile (Go is required):

	go build build.go

To run without building:

	go run build.go
*/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	Version            = "1.1.0"
	rootDirectory      = "nsroot"
	snippetOutputPath  = "snippets/snippets.json"
	documentOutputPath = "COMMANDS.md"
	tabSpace           = "  "
	snippetOutputMode  = 0o755
	documentOutputMode = 0o666
)

type Snippet struct {
	Prefix      any    `json:"prefix"`
	Body        any    `json:"body"`
	Description string `json:"description"`
}

type namedSnippet struct {
	name      string
	namespace string
	snippet   Snippet
}

func main() {
	if len(os.Args) > 1 {
		fmt.Printf(`Shellman build tool v%v

This tool doesn't accept any argument. Run it from project root directory.
It concatenates 'nsroot' snippets to 'snippets/snippets.json' and generates 'COMMANDS.md'.
`, Version)
		os.Exit(1)
	}

	if err := build("."); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func build(projectRoot string) error {
	snippetJSON, documentation, err := generate(projectRoot)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(projectRoot, snippetOutputPath), snippetJSON, snippetOutputMode); err != nil {
		return fmt.Errorf("write snippet file: %w", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, documentOutputPath), documentation, documentOutputMode); err != nil {
		return fmt.Errorf("write documentation file: %w", err)
	}
	return nil
}

func generate(projectRoot string) ([]byte, []byte, error) {
	snippets, err := readSnippets(filepath.Join(projectRoot, rootDirectory))
	if err != nil {
		return nil, nil, err
	}

	snippetJSON, err := renderSnippetJSON(snippets)
	if err != nil {
		return nil, nil, err
	}
	return snippetJSON, renderDocumentation(snippets), nil
}

func readSnippets(root string) ([]namedSnippet, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read namespace root %q: %w", root, err)
	}

	var snippets []namedSnippet
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("inspect namespace entry %q: %w", entry.Name(), err)
		}
		if !info.IsDir() {
			continue
		}

		namespace := entry.Name()
		namespacePath := filepath.Join(root, namespace)
		files, err := os.ReadDir(namespacePath)
		if err != nil {
			return nil, fmt.Errorf("read namespace %q: %w", namespace, err)
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				return nil, fmt.Errorf("inspect %q: %w", filepath.Join(namespacePath, file.Name()), err)
			}
			if fileInfo.IsDir() {
				return nil, fmt.Errorf("namespace directories should not contain nested directories: %q contains %q", namespace, file.Name())
			}

			path := filepath.Join(namespacePath, file.Name())
			snippet, err := readSnippet(path)
			if err != nil {
				return nil, err
			}
			snippets = append(snippets, namedSnippet{
				name:      namespace + "." + strings.TrimSuffix(file.Name(), ".json"),
				namespace: namespace,
				snippet:   snippet,
			})
		}
	}
	return snippets, nil
}

func readSnippet(path string) (Snippet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snippet{}, fmt.Errorf("read snippet %q: %w", path, err)
	}

	var snippet Snippet
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&snippet); err != nil {
		return Snippet{}, fmt.Errorf("decode snippet %q: %w", path, err)
	}
	if err := ensureJSONEnd(decoder); err != nil {
		return Snippet{}, fmt.Errorf("decode snippet %q: %w", path, err)
	}
	if err := validateSnippet(snippet); err != nil {
		return Snippet{}, fmt.Errorf("validate snippet %q: %w", path, err)
	}
	return snippet, nil
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var extra any
	err := decoder.Decode(&extra)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err == nil {
		return errors.New("multiple JSON values")
	}
	return err
}

func validateSnippet(snippet Snippet) error {
	switch prefix := snippet.Prefix.(type) {
	case string:
		if prefix == "" {
			return errors.New("prefix must not be empty")
		}
	case []any:
		if len(prefix) == 0 {
			return errors.New("prefix list must not be empty")
		}
		for _, value := range prefix {
			if text, ok := value.(string); !ok || text == "" {
				return errors.New("prefix list must contain non-empty strings")
			}
		}
	default:
		return errors.New("prefix must be a string or a list of strings")
	}

	switch body := snippet.Body.(type) {
	case string:
	case []any:
		for _, value := range body {
			if _, ok := value.(string); !ok {
				return errors.New("body list must contain strings")
			}
		}
	default:
		return errors.New("body must be a string or a list of strings")
	}
	return nil
}

func renderSnippetJSON(ordered []namedSnippet) ([]byte, error) {
	snippets := make(map[string]Snippet, len(ordered))
	for _, item := range ordered {
		if _, exists := snippets[item.name]; exists {
			return nil, fmt.Errorf("duplicate snippet name %q", item.name)
		}
		snippets[item.name] = item.snippet
	}

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", tabSpace)
	if err := encoder.Encode(snippets); err != nil {
		return nil, fmt.Errorf("encode snippets: %w", err)
	}
	return output.Bytes(), nil
}

func renderDocumentation(snippets []namedSnippet) []byte {
	var titles strings.Builder
	var body strings.Builder
	titles.WriteString("# Commands\n\n")

	currentNamespace := ""
	for _, item := range snippets {
		if item.namespace != currentNamespace {
			currentNamespace = item.namespace
			titles.WriteString("### " + currentNamespace + "\n\n")
		}

		prefixes := snippetPrefixes(item.snippet.Prefix)
		title := strings.Join(prefixes, " , ")
		titles.WriteString("  - [" + title + "](#" + strings.ReplaceAll(title, " ", "-") + ")\n\n")
		body.WriteString("## " + title + "\n\n")
		body.WriteString(item.snippet.Description + "[&uarr;](#" + item.namespace + ")\n\n")
	}
	return []byte(titles.String() + body.String())
}

func snippetPrefixes(prefix any) []string {
	if text, ok := prefix.(string); ok {
		return []string{text}
	}
	values := prefix.([]any)
	prefixes := make([]string, len(values))
	for index, value := range values {
		prefixes[index] = value.(string)
	}
	return prefixes
}

// sortedSnippetNames exposes the generator's stable JSON order to tests without
// coupling generation itself to map iteration order.
func sortedSnippetNames(snippets []namedSnippet) []string {
	names := make([]string, len(snippets))
	for index, item := range snippets {
		names[index] = item.name
	}
	sort.Strings(names)
	return names
}
