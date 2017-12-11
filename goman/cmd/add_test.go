package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gostores/configurator"
)

// in GOPATH, adds "test" command
// and compares the content of all files in cmd directory of testproject
// with appropriate golden files.
// Use -update to update existing golden files.
func TestGoldenAddCmd(t *testing.T) {
	projectName := "github.com/gostores/testproject"
	project := NewProject(projectName)
	defer os.RemoveAll(project.AbsPath())

	configurator.Set("author", "NAME HERE <EMAIL ADDRESS>")
	configurator.Set("license", "apache")
	configurator.Set("year", 2017)
	defer configurator.Set("author", nil)
	defer configurator.Set("license", nil)
	defer configurator.Set("year", nil)

	// Initialize the project first.
	initializeProject(project)

	// Then add the "test" command.
	cmdName := "test"
	cmdPath := filepath.Join(project.CmdPath(), cmdName+".go")
	createCmdFile(project.License(), cmdPath, cmdName)

	expectedFiles := []string{".", "root.go", "test.go"}
	gotFiles := []string{}

	// Check project file hierarchy and compare the content of every single file
	// with appropriate golden file.
	err := filepath.Walk(project.CmdPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Make path relative to project.CmdPath().
		// then it returns just "root.go".
		relPath, err := filepath.Rel(project.CmdPath(), path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		gotFiles = append(gotFiles, relPath)
		goldenPath := filepath.Join("testdata", filepath.Base(path)+".golden")

		switch relPath {
		// Known directories.
		case ".":
			return nil
		// Known files.
		case "root.go", "test.go":
			if *update {
				got, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				ioutil.WriteFile(goldenPath, got, 0644)
			}
			return compareFiles(path, goldenPath)
		}
		// Unknown file.
		return errors.New("unknown file: " + path)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check if some files lack.
	if err := checkLackFiles(expectedFiles, gotFiles); err != nil {
		t.Fatal(err)
	}
}

func TestValidateCmdName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"cmdName", "cmdName"},
		{"cmd_name", "cmdName"},
		{"cmd-name", "cmdName"},
		{"cmd______Name", "cmdName"},
		{"cmd------Name", "cmdName"},
		{"cmd______name", "cmdName"},
		{"cmd------name", "cmdName"},
		{"cmdName-----", "cmdName"},
		{"cmdname-", "cmdname"},
	}

	for _, testCase := range testCases {
		got := validateCmdName(testCase.input)
		if testCase.expected != got {
			t.Errorf("Expected %q, got %q", testCase.expected, got)
		}
	}
}
