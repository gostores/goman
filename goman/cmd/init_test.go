package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/govenue/configurator"
)

func TestGoldenInitCmd(t *testing.T) {
	projectName := "github.com/govenue/testproject"
	project := NewProject(projectName)
	defer os.RemoveAll(project.AbsPath())

	configurator.Set("author", "NAME HERE <EMAIL ADDRESS>")
	configurator.Set("license", "apache")
	configurator.Set("year", 2017)
	defer configurator.Set("author", nil)
	defer configurator.Set("license", nil)
	defer configurator.Set("year", nil)

	os.Args = []string{"goman", "init", projectName}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal("Error by execution:", err)
	}

	expectedFiles := []string{".", "cmd", "LICENSE", "main.go", "cmd/root.go"}
	gotFiles := []string{}

	// Check project file hierarchy and compare the content of every single file
	// with appropriate golden file.
	err := filepath.Walk(project.AbsPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Make path relative to project.AbsPath().
		// then it returns just "cmd/root.go".
		relPath, err := filepath.Rel(project.AbsPath(), path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		gotFiles = append(gotFiles, relPath)
		goldenPath := filepath.Join("testdata", filepath.Base(path)+".golden")

		switch relPath {
		// Known directories.
		case ".", "cmd":
			return nil
		// Known files.
		case "LICENSE", "main.go", "cmd/root.go":
			if *update {
				got, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				if err := ioutil.WriteFile(goldenPath, got, 0644); err != nil {
					t.Fatal("Error while updating file:", err)
				}
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
