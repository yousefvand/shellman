package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const (
	historicalV6SnippetSHA256    = "5840e7cc9ff1e1ba17413f2c54f242224aac413a94afb89a77cc6cf89c668d4b"
	historicalV6CommandsSHA256   = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration1SHA256   = "42f7919542f6ffb7d48fb8e44fd3a4091eb4a63f3d9774b594d526f56b6408b1"
	approvedV7Iteration1Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration2SHA256   = "7ddfd8260bdc5bddd951abe4712582b6e164d6ecdab693fcee7275ae526827b8"
	approvedV7Iteration2Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration3SHA256   = "04a2f6f0d046a828ac0d48d94e08469184232879e75859140ab28fbb28130902"
	approvedV7Iteration3Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration4SHA256   = "5929c529f11a3e9c01c7d80289025873ce9b8613e7db3d20606d0d8f70c15e32"
	approvedV7Iteration4Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration5SHA256   = "3310be931efc6aba4a25804391a5011c956f698ce0818d2a36d115bcb1f930c7"
	approvedV7Iteration5Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration6SHA256   = "fdbb2248248bd8dadf5108b82a8521b34bb7ae7c0fcffc0a81e2ec18a5f7fbf9"
	approvedV7Iteration6Commands = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedMigrationName        = "archive.compress-tar-gz"
	iteration2MigrationName      = "archive.compress-tar-xz"
	iteration3MigrationName      = "archive.compress-zip"
	iteration4MigrationName      = "archive.decompress-tar-gz"
	iteration5MigrationName      = "archive.decompress-tar-xz"
	iteration6MigrationName      = "archive.decompress-unzip"
	v6CompressTarGzBody          = "tar -czvf ${1|/path/to/archive, \"${pathToArchive}\"|}.tar.gz ${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\n"
	currentCompressTarGzBody     = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\"\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\ntar -cf \"${archive_path}.tar\" \"${source_path}\" && gzip -f \"${archive_path}.tar\"\n"
	v6CompressTarXzBody          = "tar -cJf ${1|/path/to/archive, \"${pathToArchive}\"|}.tar.xz ${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\n"
	currentCompressTarXzBody     = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\"\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\ntar -cf \"${archive_path}.tar\" \"${source_path}\" && xz -f \"${archive_path}.tar\"\n"
	v6CompressZipBody            = "zip -rq ${1|/path/to/archive, \"${pathToArchive}\"|}.zip ${2|/path/to/directory-or-file,\"${pathToDirectoryOrFile}\"|}\n"
	currentCompressZipBody       = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file,\"${pathToDirectoryOrFile}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\nzip -rq \"${archive_path}.zip\" \"${source_path}\"\n"
	v6DecompressTarGzBody        = "tar -C ${1|/extract/to/path, \"${extractToPath}\"|} -xzvf ${2|/path/to/archive, \"${pathToArchive}\"|}.tar.gz\n"
	currentDecompressTarGzBody   = "extract_path=\"${1|/extract/to/path, \"${extractToPath}\"|}\"\narchive_path=\"${2|/path/to/archive, \"${pathToArchive}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\n(\n  temporary_directory=$(mktemp -d \"${TMPDIR:-/tmp}/shellman.XXXXXXXXXX\") || exit\n  status=0\n  trap 'status=$?; trap - 0; rm -rf \"${temporary_directory}\"; exit \"${status}\"' 0\n  trap 'exit 129' HUP\n  trap 'exit 130' INT\n  trap 'exit 143' TERM\n  gzip -dc \"${archive_path}.tar.gz\" > \"${temporary_directory}/archive.tar\" &&\n    (cd \"${extract_path}\" && tar -xf \"${temporary_directory}/archive.tar\")\n)\n"
	v6DecompressTarXzBody        = "tar -C ${1|/extract/to/path, \"${extractToPath}\"|} -xf ${2|/path/to/archive, \"${pathToArchive}\"|}.tar.xz\n"
	currentDecompressTarXzBody   = "extract_path=\"${1|/extract/to/path, \"${extractToPath}\"|}\"\narchive_path=\"${2|/path/to/archive, \"${pathToArchive}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\n(\n  temporary_directory=$(mktemp -d \"${TMPDIR:-/tmp}/shellman.XXXXXXXXXX\") || exit\n  status=0\n  trap 'status=$?; trap - 0; rm -rf \"${temporary_directory}\"; exit \"${status}\"' 0\n  trap 'exit 129' HUP\n  trap 'exit 130' INT\n  trap 'exit 143' TERM\n  xz -dc \"${archive_path}.tar.xz\" > \"${temporary_directory}/archive.tar\" &&\n    (cd \"${extract_path}\" && tar -xf \"${temporary_directory}/archive.tar\")\n)\n"
	v6DecompressUnzipBody        = "unzip -q ${1|/path/to/archive, \"${pathToArchive}\"|}.zip -d ${2|/extract/to/path,\"${extractToPath}\"|}"
	currentDecompressUnzipBody   = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nextract_path=\"${2|/extract/to/path,\"${extractToPath}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\nunzip -q \"${archive_path}.zip\" -d \"${extract_path}\""
)

var v6Namespaces = []string{
	"archive", "array", "command", "cryptography", "date", "event",
	"filesystem", "float", "fn-fx", "ftp", "function", "git", "http",
	"input", "integer", "internal", "ip", "math", "misc", "output",
	"process", "string", "system", "time", "variable",
}

func TestGeneratedOutputMatchesFiles(t *testing.T) {
	snippetJSON, commands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}

	assertFileBytes(t, snippetOutputPath, snippetJSON)
	assertFileBytes(t, documentOutputPath, commands)
	if !bytes.HasSuffix(snippetJSON, []byte("\n")) {
		t.Error("snippet output must end with a newline")
	}
	if !bytes.HasSuffix(commands, []byte("\n\n")) {
		t.Error("documentation output must end with two newlines")
	}
}

func TestMigrationChangesOnlyApprovedSnippet(t *testing.T) {
	ordered, err := readSnippets(rootDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ordered), 278; got != want {
		t.Fatalf("snippet count = %d, want %d", got, want)
	}
	if got := ordered[3].name; got != iteration4MigrationName {
		t.Fatalf("fourth snippet in actual nsroot traversal = %q, want %q", got, iteration4MigrationName)
	}
	if got := ordered[4].name; got != iteration5MigrationName {
		t.Fatalf("fifth snippet in actual nsroot traversal = %q, want %q", got, iteration5MigrationName)
	}
	if got := ordered[5].name; got != iteration6MigrationName {
		t.Fatalf("sixth snippet in actual nsroot traversal = %q, want %q", got, iteration6MigrationName)
	}
	currentIteration6, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "approved v7 iteration-6 snippets (six changed, 272 unchanged)", currentIteration6, approvedV7Iteration6SHA256)
	assertSHA256(t, "approved v7 iteration-6 commands", renderDocumentation(ordered), approvedV7Iteration6Commands)

	foundFirst := false
	foundSecond := false
	foundThird := false
	foundFourth := false
	foundFifth := false
	foundSixth := false
	for index := range ordered {
		switch ordered[index].name {
		case approvedMigrationName:
			foundFirst = true
			wantPrefix := []any{"archive compress tar.gz", "archive tar.gz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("approved snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .tar.gz file"; got != want {
				t.Fatalf("approved snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressTarGzBody {
				t.Fatalf("approved snippet body = %#v, want %#v", got, currentCompressTarGzBody)
			}
		case iteration2MigrationName:
			foundSecond = true
			wantPrefix := []any{"archive compress tar.xz", "archive tar.xz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-2 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .tar.xz file"; got != want {
				t.Fatalf("iteration-2 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressTarXzBody {
				t.Fatalf("iteration-2 snippet body = %#v, want %#v", got, currentCompressTarXzBody)
			}
		case iteration3MigrationName:
			foundThird = true
			wantPrefix := []any{"archive compress .zip", "archive zip"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-3 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .zip file"; got != want {
				t.Fatalf("iteration-3 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressZipBody {
				t.Fatalf("iteration-3 snippet body = %#v, want %#v", got, currentCompressZipBody)
			}
		case iteration4MigrationName:
			foundFourth = true
			wantPrefix := []any{"archive decompress tar.gz", "decompress tar.gz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-4 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .tar.gz file to specified path"; got != want {
				t.Fatalf("iteration-4 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressTarGzBody {
				t.Fatalf("iteration-4 snippet body = %#v, want %#v", got, currentDecompressTarGzBody)
			}
		case iteration5MigrationName:
			foundFifth = true
			wantPrefix := []any{"archive decompress tar.xz", "decompress tar.xz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-5 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .tar.xz file to specified path"; got != want {
				t.Fatalf("iteration-5 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressTarXzBody {
				t.Fatalf("iteration-5 snippet body = %#v, want %#v", got, currentDecompressTarXzBody)
			}
		case iteration6MigrationName:
			foundSixth = true
			wantPrefix := []any{"archive decompress .zip", "archive unzip"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-6 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .zip file to specified path"; got != want {
				t.Fatalf("iteration-6 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressUnzipBody {
				t.Fatalf("iteration-6 snippet body = %#v, want %#v", got, currentDecompressUnzipBody)
			}
			ordered[index].snippet.Body = v6DecompressUnzipBody
		}
	}
	if !foundFirst {
		t.Fatalf("approved snippet %q is missing", approvedMigrationName)
	}
	if !foundSecond {
		t.Fatalf("iteration-2 snippet %q is missing", iteration2MigrationName)
	}
	if !foundThird {
		t.Fatalf("iteration-3 snippet %q is missing", iteration3MigrationName)
	}
	if !foundFourth {
		t.Fatalf("iteration-4 snippet %q is missing", iteration4MigrationName)
	}
	if !foundFifth {
		t.Fatalf("iteration-5 snippet %q is missing", iteration5MigrationName)
	}
	if !foundSixth {
		t.Fatalf("iteration-6 snippet %q is missing", iteration6MigrationName)
	}

	reconstructedIteration5, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-5 snippets", reconstructedIteration5, approvedV7Iteration5SHA256)
	assertSHA256(t, "approved v7 iteration-5 commands", renderDocumentation(ordered), approvedV7Iteration5Commands)

	for index := range ordered {
		if ordered[index].name == iteration5MigrationName {
			ordered[index].snippet.Body = v6DecompressTarXzBody
		}
	}

	reconstructedIteration4, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-4 snippets", reconstructedIteration4, approvedV7Iteration4SHA256)
	assertSHA256(t, "approved v7 iteration-4 commands", renderDocumentation(ordered), approvedV7Iteration4Commands)

	for index := range ordered {
		if ordered[index].name == iteration4MigrationName {
			ordered[index].snippet.Body = v6DecompressTarGzBody
		}
	}

	reconstructedIteration3, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-3 snippets", reconstructedIteration3, approvedV7Iteration3SHA256)
	assertSHA256(t, "approved v7 iteration-3 commands", renderDocumentation(ordered), approvedV7Iteration3Commands)

	for index := range ordered {
		if ordered[index].name == iteration3MigrationName {
			ordered[index].snippet.Body = v6CompressZipBody
		}
	}

	reconstructedIteration2, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-2 snippets", reconstructedIteration2, approvedV7Iteration2SHA256)
	assertSHA256(t, "approved v7 iteration-2 commands", renderDocumentation(ordered), approvedV7Iteration2Commands)

	for index := range ordered {
		if ordered[index].name == iteration2MigrationName {
			ordered[index].snippet.Body = v6CompressTarXzBody
		}
	}

	reconstructedIteration1, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-1 snippets", reconstructedIteration1, approvedV7Iteration1SHA256)
	assertSHA256(t, "approved v7 iteration-1 commands", renderDocumentation(ordered), approvedV7Iteration1Commands)

	for index := range ordered {
		if ordered[index].name == approvedMigrationName {
			ordered[index].snippet.Body = v6CompressTarGzBody
		}
	}
	reconstructedV6, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed historical v6 snippets", reconstructedV6, historicalV6SnippetSHA256)
	assertSHA256(t, "historical v6 commands", renderDocumentation(ordered), historicalV6CommandsSHA256)
}

func TestCompressTarGzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentCompressTarGzBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -cf "${archive_path}.tar" "${source_path}"`) {
		t.Fatal("snippet must use the portable tar -c and -f options")
	}
	for _, nonPortable := range []string{" -z", " --", " -C"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestCompressTarGzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "gzip")

	script := runnableGeneratedCompressTarGz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".tar.gz"
		run(t, root, "gzip", "-t", archive)
		listing := run(t, root, "tar", "-tzf", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, extract, "tar", "-xzf", archive)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptPath(t, script), archiveBase, "empty source")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without gzip output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".tar.gz"); !os.IsNotExist(err) {
			t.Fatalf("gzip archive exists after tar failure: %v", err)
		}
	})

	t.Run("existing destination is overwritten", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement\n"))
		archiveBase := filepath.Join(root, "existing")
		mustWriteFile(t, archiveBase+".tar.gz", []byte("old invalid archive"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		run(t, root, "gzip", "-t", archiveBase+".tar.gz")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if listing != "source.txt\n" {
			t.Fatalf("replacement archive listing = %q", listing)
		}
	})

	t.Run("tar failure prevents gzip", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 19\n")
		marker := filepath.Join(root, "gzip-ran")
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\n: > \"$SHELLMAN_GZIP_MARKER\"\nexit 0\n")
		environment := append(os.Environ(), "PATH="+fakeBin, "SHELLMAN_GZIP_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 19 {
			t.Fatalf("exit status = %d, want tar status 19: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("gzip ran after tar failure: %v", err)
		}
	})

	t.Run("gzip failure leaves intermediate tar", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nexit 23\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		archiveBase := filepath.Join(root, "gzip failure")
		_, err := runCommand(root, environment, "sh", scriptFile, archiveBase, "source.txt")
		if exitCode(err) != 23 {
			t.Fatalf("exit status = %d, want gzip status 23: %v", exitCode(err), err)
		}
		if _, err := os.Stat(archiveBase + ".tar"); err != nil {
			t.Fatalf("tar intermediate not retained after gzip failure: %v", err)
		}
		if _, err := os.Stat(archiveBase + ".tar.gz"); !os.IsNotExist(err) {
			t.Fatalf("gzip output exists after simulated failure: %v", err)
		}
	})

	t.Run("source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "hyphen archive")
		run(t, root, "sh", scriptFile, archiveBase, "-source")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if !strings.Contains(listing, "./-source/file.txt\n") {
			t.Fatalf("hyphen-source archive listing = %q", listing)
		}
	})

	t.Run("destination inside source directory", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source"
		mustMkdirAll(t, filepath.Join(root, sourceName))
		mustWriteFile(t, filepath.Join(root, sourceName, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, sourceName, "inside")
		run(t, root, "sh", scriptFile, archiveBase, sourceName)
		archive := archiveBase + ".tar.gz"
		run(t, root, "gzip", "-t", archive)
		listing := run(t, root, "tar", "-tzf", archive)
		if !strings.Contains(listing, "source/file.txt\n") {
			t.Fatalf("inside-source archive lacks input file: %q", listing)
		}
		if strings.Contains(listing, "source/inside.tar") || strings.Contains(listing, "source/inside.tar.gz") {
			t.Fatalf("archive included itself: %q", listing)
		}
	})
}

func runnableGeneratedCompressTarGz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[approvedMigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", approvedMigrationName)
	}
	if body != currentCompressTarGzBody {
		t.Fatal("generated snippet body differs from the approved body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestCompressTarXzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentCompressTarXzBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -cf "${archive_path}.tar" "${source_path}"`) {
		t.Fatal("snippet must use the portable tar -c and -f options")
	}
	for _, nonPortable := range []string{" -J", " --", " -C"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestCompressTarXzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "xz")

	script := runnableGeneratedCompressTarXz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman xz archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".tar.xz"
		run(t, root, "xz", "-t", archive)
		listing := run(t, root, "tar", "-tJf", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, extract, "tar", "-xJf", archive)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman xz archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptPath(t, script), archiveBase, "empty source")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without xz output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".tar.xz"); !os.IsNotExist(err) {
			t.Fatalf("xz archive exists after tar failure: %v", err)
		}
	})

	t.Run("existing destination is overwritten", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement\n"))
		archiveBase := filepath.Join(root, "existing")
		mustWriteFile(t, archiveBase+".tar.xz", []byte("old invalid archive"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		run(t, root, "xz", "-t", archiveBase+".tar.xz")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if listing != "source.txt\n" {
			t.Fatalf("replacement archive listing = %q", listing)
		}
	})

	t.Run("tar failure prevents xz", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 29\n")
		marker := filepath.Join(root, "xz-ran")
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\n: > \"$SHELLMAN_XZ_MARKER\"\nexit 0\n")
		environment := append(os.Environ(), "PATH="+fakeBin, "SHELLMAN_XZ_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 29 {
			t.Fatalf("exit status = %d, want tar status 29: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("xz ran after tar failure: %v", err)
		}
	})

	t.Run("xz failure leaves intermediate tar", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nexit 31\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		archiveBase := filepath.Join(root, "xz failure")
		_, err := runCommand(root, environment, "sh", scriptFile, archiveBase, "source.txt")
		if exitCode(err) != 31 {
			t.Fatalf("exit status = %d, want xz status 31: %v", exitCode(err), err)
		}
		if _, err := os.Stat(archiveBase + ".tar"); err != nil {
			t.Fatalf("tar intermediate not retained after xz failure: %v", err)
		}
		if _, err := os.Stat(archiveBase + ".tar.xz"); !os.IsNotExist(err) {
			t.Fatalf("xz output exists after simulated failure: %v", err)
		}
	})

	t.Run("source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "hyphen archive")
		run(t, root, "sh", scriptFile, archiveBase, "-source")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if !strings.Contains(listing, "./-source/file.txt\n") {
			t.Fatalf("hyphen-source archive listing = %q", listing)
		}
	})

	t.Run("destination inside source directory", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source"
		mustMkdirAll(t, filepath.Join(root, sourceName))
		mustWriteFile(t, filepath.Join(root, sourceName, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, sourceName, "inside")
		run(t, root, "sh", scriptFile, archiveBase, sourceName)
		archive := archiveBase + ".tar.xz"
		run(t, root, "xz", "-t", archive)
		listing := run(t, root, "tar", "-tJf", archive)
		if !strings.Contains(listing, "source/file.txt\n") {
			t.Fatalf("inside-source archive lacks input file: %q", listing)
		}
		if strings.Contains(listing, "source/inside.tar") || strings.Contains(listing, "source/inside.tar.xz") {
			t.Fatalf("archive included itself: %q", listing)
		}
	})
}

func runnableGeneratedCompressTarXz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration2MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration2MigrationName)
	}
	if body != currentCompressTarXzBody {
		t.Fatal("generated iteration-2 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestCompressZipPlaceholderContractAndOptionSafety(t *testing.T) {
	body := currentCompressZipBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file,"${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `zip -rq "${archive_path}.zip" "${source_path}"`) {
		t.Fatal("snippet must quote the archive and source paths")
	}
	if strings.Contains(body, " --") {
		t.Fatal("snippet must not rely on non-portable end-of-options syntax")
	}
}

func TestCompressZipBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "zip")
	requireCommand(t, "unzip")

	script := runnableGeneratedCompressZip(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman zip archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".zip"
		run(t, root, "unzip", "-t", archive)
		listing := run(t, root, "unzip", "-Z1", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, root, "unzip", "-q", archive, "-d", extract)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman zip archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptFile, archiveBase, "empty source")
		listing := run(t, root, "unzip", "-Z1", archiveBase+".zip")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".zip"); !os.IsNotExist(err) {
			t.Fatalf("archive exists after failure: %v", err)
		}
	})

	t.Run("existing archive is updated", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("first\n"))
		archiveBase := filepath.Join(root, "existing")
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement content\n"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		got := run(t, root, "unzip", "-p", archiveBase+".zip", "source.txt")
		if got != "replacement content\n" {
			t.Fatalf("updated content = %q", got)
		}
	})

	t.Run("zip failure status propagates", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "zip"), "#!/bin/sh\nexit 37\n")
		environment := append(os.Environ(), "PATH="+fakeBin)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 37 {
			t.Fatalf("exit status = %d, want zip status 37: %v", exitCode(err), err)
		}
	})

	t.Run("archive and source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		run(t, root, "sh", scriptFile, "-archive", "-source")
		listing := run(t, root, "unzip", "-Z1", filepath.Join(root, "-archive.zip"))
		if !strings.Contains(listing, "-source/file.txt\n") {
			t.Fatalf("hyphen-path archive listing = %q", listing)
		}
	})

	t.Run("missing destination directory fails", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "missing", "archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "source.txt")
		if err == nil {
			t.Fatal("missing destination directory unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".zip"); !os.IsNotExist(err) {
			t.Fatalf("archive exists after destination failure: %v", err)
		}
	})
}

func runnableGeneratedCompressZip(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration3MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration3MigrationName)
	}
	if body != currentCompressZipBody {
		t.Fatal("generated iteration-3 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file,"${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestDecompressTarGzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentDecompressTarGzBody
	for _, placeholder := range []string{
		`${1|/extract/to/path, "${extractToPath}"|}`,
		`${2|/path/to/archive, "${pathToArchive}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -xf "${temporary_directory}/archive.tar"`) {
		t.Fatal("snippet must use portable tar -x and -f options")
	}
	for _, nonPortable := range []string{" -z", " -C", " --"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
	for _, trap := range []string{
		`trap 'status=$?; trap - 0; rm -rf "${temporary_directory}"; exit "${status}"' 0`,
		`trap 'exit 129' HUP`,
		`trap 'exit 130' INT`,
		`trap 'exit 143' TERM`,
	} {
		if !strings.Contains(body, trap) {
			t.Fatalf("snippet lacks separate cleanup/signal trap %q", trap)
		}
	}
}

func TestDecompressTarGzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "gzip")
	requireCommand(t, "mktemp")

	script := runnableGeneratedDecompressTarGz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters integrity and cleanup", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		sourceName := "source dir [brackets] #dollar$"
		mustMkdirAll(t, filepath.Join(staging, sourceName, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, sourceName, "file name [1] #$.txt"), []byte("decompress integrity\n"))
		archiveBase := filepath.Join(root, "archive [1] #$")
		createTarGz(t, staging, archiveBase, sourceName)
		extract := filepath.Join(root, "extract dir [1] #$")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "temporary root")
		mustMkdirAll(t, temporaryRoot)

		runWithEnvironment(t, root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, archiveBase)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "decompress integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarGz(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		run(t, root, "sh", scriptFile, extract, "-archive")
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("missing archive propagates gzip failure and cleans up", func(t *testing.T) {
		root := t.TempDir()
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		_, err := runCommand(root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, filepath.Join(root, "missing"))
		if err == nil {
			t.Fatal("missing archive unexpectedly succeeded")
		}
		assertDirectoryEmpty(t, extract)
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("gzip failure prevents tar and preserves status", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nexit 41\n")
		marker := filepath.Join(root, "tar-ran")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\nexit 0\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 41 {
			t.Fatalf("exit status = %d, want gzip status 41: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("tar ran after gzip failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	for _, test := range []struct {
		name       string
		signalFlag string
		exitStatus int
	}{
		{name: "INT cleans up and exits 130", signalFlag: "-INT", exitStatus: 130},
		{name: "TERM cleans up and exits 143", signalFlag: "-TERM", exitStatus: 143},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			fakeBin := filepath.Join(root, "bin")
			mustMkdirAll(t, fakeBin)
			mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nkill "+test.signalFlag+" \"$PPID\"\nexit 0\n")
			marker := filepath.Join(root, "tar-ran")
			mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\nexit 0\n")
			extract := filepath.Join(root, "extract")
			mustMkdirAll(t, extract)
			temporaryRoot := filepath.Join(root, "tmp")
			mustMkdirAll(t, temporaryRoot)
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
			_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
			if exitCode(err) != test.exitStatus {
				t.Fatalf("exit status = %d, want %d: %v", exitCode(err), test.exitStatus, err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("tar ran after %s: %v", test.signalFlag, err)
			}
			assertDirectoryEmpty(t, temporaryRoot)
		})
	}

	t.Run("tar failure propagates and cleans up", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarGz(t, root, filepath.Join(root, "archive"), "file.txt")
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 43\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 43 {
			t.Fatalf("exit status = %d, want tar status 43: %v", exitCode(err), err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("mktemp failure prevents decompression", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "mktemp"), "#!/bin/sh\nexit 47\n")
		marker := filepath.Join(root, "gzip-ran")
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\n: > \"$SHELLMAN_GZIP_MARKER\"\nexit 0\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_GZIP_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 47 {
			t.Fatalf("exit status = %d, want mktemp status 47: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("gzip ran after mktemp failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("missing extraction directory fails after decompression", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "archive")
		createTarGz(t, root, archiveBase, "file.txt")
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		_, err := runCommand(root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, filepath.Join(root, "missing"), archiveBase)
		if err == nil {
			t.Fatal("missing extraction directory unexpectedly succeeded")
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})
}

func runnableGeneratedDecompressTarGz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration4MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration4MigrationName)
	}
	if body != currentDecompressTarGzBody {
		t.Fatal("generated iteration-4 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/extract/to/path, "${extractToPath}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/archive, "${pathToArchive}"|}`, `$2`)
}

func TestDecompressTarXzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentDecompressTarXzBody
	for _, placeholder := range []string{
		`${1|/extract/to/path, "${extractToPath}"|}`,
		`${2|/path/to/archive, "${pathToArchive}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	for _, required := range []string{`xz -dc "${archive_path}.tar.xz"`, `tar -xf "${temporary_directory}/archive.tar"`, `trap 'exit 130' INT`, `trap 'exit 143' TERM`} {
		if !strings.Contains(body, required) {
			t.Fatalf("snippet lacks %q", required)
		}
	}
	for _, nonPortable := range []string{" -J", " -C", " --"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestDecompressTarXzBehavior(t *testing.T) {
	for _, command := range []string{"sh", "tar", "xz", "mktemp"} {
		requireCommand(t, command)
	}
	scriptFile := scriptPath(t, runnableGeneratedDecompressTarXz(t))
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("integrity special paths and normal cleanup", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		source := "source dir [5] #$"
		mustMkdirAll(t, filepath.Join(staging, source, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, source, "file [5] #$.txt"), []byte("iteration five\n"))
		archiveBase := filepath.Join(root, "archive [5] #$")
		createTarXz(t, staging, archiveBase, source)
		extract := filepath.Join(root, "extract [5] #$")
		temporaryRoot := filepath.Join(root, "tmp root")
		mustMkdirAll(t, extract)
		mustMkdirAll(t, temporaryRoot)
		runWithEnvironment(t, root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, archiveBase)
		got, err := os.ReadFile(filepath.Join(extract, source, "file [5] #$.txt"))
		if err != nil || string(got) != "iteration five\n" {
			t.Fatalf("extracted content = %q, %v", got, err)
		}
		if info, err := os.Stat(filepath.Join(extract, source, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarXz(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		run(t, root, "sh", scriptFile, extract, "-archive")
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("xz failure prevents tar and cleans up", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nexit 51\n")
		marker := filepath.Join(root, "tar-ran")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 51 {
			t.Fatalf("exit status = %d, want 51: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("tar ran after xz failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("tar failure propagates and cleans up", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nprintf archive\n")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 53\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 53 {
			t.Fatalf("exit status = %d, want 53: %v", exitCode(err), err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("mktemp failure prevents xz", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "mktemp"), "#!/bin/sh\nexit 55\n")
		marker := filepath.Join(root, "xz-ran")
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\n: > \"$SHELLMAN_XZ_MARKER\"\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_XZ_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 55 {
			t.Fatalf("exit status = %d, want 55: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("xz ran after mktemp failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	for _, test := range []struct {
		name, signal string
		status       int
	}{{"INT cleanup", "-INT", 130}, {"TERM cleanup", "-TERM", 143}} {
		t.Run(test.name, func(t *testing.T) {
			root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
			mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nkill "+test.signal+" \"$PPID\"\n")
			_, err := runCommand(root, append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, filepath.Join(root, "archive"))
			if exitCode(err) != test.status {
				t.Fatalf("exit status = %d, want %d: %v", exitCode(err), test.status, err)
			}
			assertDirectoryEmpty(t, temporaryRoot)
		})
	}
}

func runnableGeneratedDecompressTarXz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration5MigrationName].Body.(string)
	if !ok || body != currentDecompressTarXzBody {
		t.Fatal("generated iteration-5 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1|/extract/to/path, "${extractToPath}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/archive, "${pathToArchive}"|}`, `$2`)
}

func decompressionFailurePaths(t *testing.T) (root, fakeBin, extract, temporaryRoot string) {
	t.Helper()
	root = t.TempDir()
	fakeBin = filepath.Join(root, "bin")
	extract = filepath.Join(root, "extract")
	temporaryRoot = filepath.Join(root, "tmp")
	mustMkdirAll(t, fakeBin)
	mustMkdirAll(t, extract)
	mustMkdirAll(t, temporaryRoot)
	return
}

func createTarXz(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "tar", "-cf", archiveBase+".tar", source)
	run(t, dir, "xz", "-f", archiveBase+".tar")
}

func TestDecompressUnzipPlaceholderContractAndOptionSafety(t *testing.T) {
	body := currentDecompressUnzipBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/extract/to/path,"${extractToPath}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `unzip -q "${archive_path}.zip" -d "${extract_path}"`) {
		t.Fatal("snippet must quote archive and extraction paths")
	}
	if strings.Contains(body, " --") {
		t.Fatal("snippet must not rely on non-portable end-of-options syntax")
	}
}

func TestDecompressUnzipBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "zip")
	requireCommand(t, "unzip")
	scriptFile := scriptPath(t, runnableGeneratedDecompressUnzip(t))
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("integrity spaces special characters and empty directory", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		source := "source dir [6] #$"
		mustMkdirAll(t, filepath.Join(staging, source, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, source, "file [6] #$.txt"), []byte("iteration six\n"))
		archiveBase := filepath.Join(root, "archive [6] #$")
		createZip(t, staging, archiveBase, source)
		extract := filepath.Join(root, "extract dir [6] #$")
		run(t, root, "sh", scriptFile, archiveBase, extract)
		got, err := os.ReadFile(filepath.Join(extract, source, "file [6] #$.txt"))
		if err != nil || string(got) != "iteration six\n" {
			t.Fatalf("extracted content = %q, %v", got, err)
		}
		if info, err := os.Stat(filepath.Join(extract, source, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createZip(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		run(t, root, "sh", scriptFile, "-archive", extract)
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("missing archive fails", func(t *testing.T) {
		root := t.TempDir()
		extract := filepath.Join(root, "extract")
		_, err := runCommand(root, nil, "sh", scriptFile, filepath.Join(root, "missing"), extract)
		if err == nil {
			t.Fatal("missing archive unexpectedly succeeded")
		}
		if _, err := os.Stat(extract); !os.IsNotExist(err) {
			t.Fatalf("extraction directory exists after missing archive: %v", err)
		}
	})

	t.Run("unzip failure status propagates", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "unzip"), "#!/bin/sh\nexit 61\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), filepath.Join(root, "extract"))
		if exitCode(err) != 61 {
			t.Fatalf("exit status = %d, want 61: %v", exitCode(err), err)
		}
	})

	t.Run("unwritable extraction target fails", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "archive")
		createZip(t, root, archiveBase, "file.txt")
		notDirectory := filepath.Join(root, "not-directory")
		mustWriteFile(t, notDirectory, []byte("occupied\n"))
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, notDirectory)
		if err == nil {
			t.Fatal("file extraction target unexpectedly succeeded")
		}
	})
}

func runnableGeneratedDecompressUnzip(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration6MigrationName].Body.(string)
	if !ok || body != currentDecompressUnzipBody {
		t.Fatal("generated iteration-6 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/extract/to/path,"${extractToPath}"|}`, `$2`)
}

func createZip(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "zip", "-rq", archiveBase+".zip", source)
}

func createTarGz(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "tar", "-cf", archiveBase+".tar", source)
	run(t, dir, "gzip", "-f", archiveBase+".tar")
}

func runWithEnvironment(t *testing.T, dir string, environment []string, name string, args ...string) string {
	t.Helper()
	output, err := runCommand(dir, environment, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return string(output)
}

func assertDirectoryEmpty(t *testing.T, path string) {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("directory %s is not empty: %v", path, entries)
	}
}

func requireCommand(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Fatalf("required command %q is unavailable: %v", name, err)
	}
}

func scriptPath(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "compress-tar-gz.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	output, err := runCommand(dir, nil, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return string(output)
}

func runCommand(dir string, environment []string, name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	command.Dir = dir
	if environment != nil {
		command.Env = environment
	}
	return command.CombinedOutput()
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return -1
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestV6SnippetInventoryAndOrder(t *testing.T) {
	ordered, err := readSnippets(rootDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ordered), 278; got != want {
		t.Fatalf("snippet count = %d, want %d", got, want)
	}

	var namespaces []string
	for _, item := range ordered {
		if len(namespaces) == 0 || namespaces[len(namespaces)-1] != item.namespace {
			namespaces = append(namespaces, item.namespace)
		}
	}
	if !reflect.DeepEqual(namespaces, v6Namespaces) {
		t.Fatalf("namespaces = %q, want %q", namespaces, v6Namespaces)
	}

	names := sortedSnippetNames(ordered)
	if !sort.StringsAreSorted(names) {
		t.Fatal("generated snippet names are not sorted")
	}

	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	if got := jsonObjectKeys(t, generated); !reflect.DeepEqual(got, names) {
		t.Fatal("JSON key order differs from the v6 snippet-name order")
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	firstSnippets, firstCommands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	secondSnippets, secondCommands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstSnippets, secondSnippets) || !bytes.Equal(firstCommands, secondCommands) {
		t.Fatal("repeated generation produced different output")
	}
}

func assertFileBytes(t *testing.T, path string, generated []byte) {
	t.Helper()
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(generated, want) {
		t.Fatalf("generated output differs from %s", path)
	}
}

func assertSHA256(t *testing.T, name string, data []byte, want string) {
	t.Helper()
	digest := sha256.Sum256(data)
	if got := hex.EncodeToString(digest[:]); got != want {
		t.Fatalf("%s SHA-256 = %s, want approved contract %s", name, got, want)
	}
}

func jsonObjectKeys(t *testing.T, data []byte) []string {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		t.Fatalf("decode generated JSON object: %v", err)
	}

	var keys []string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, token.(string))
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := decoder.Token(); err != nil {
		t.Fatal(err)
	}
	return keys
}

func TestGeneratedFileStructure(t *testing.T) {
	if snippetOutputPath != "snippets/snippets.json" {
		t.Fatalf("unexpected snippet output path %q", snippetOutputPath)
	}
	if documentOutputPath != "COMMANDS.md" {
		t.Fatalf("unexpected documentation output path %q", documentOutputPath)
	}
	if strings.ContainsAny(snippetOutputPath, "\\") {
		t.Fatalf("snippet output path is not portable: %q", snippetOutputPath)
	}
}
