package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	is := assert.New(t)
	tempDir, err := ioutil.TempDir("", "duffle-import-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tests := map[string]struct {
		dest        string
		fileCreated func(string) bool
	}{
		"no --destination": {
			dest: "",
			fileCreated: func(bundleDir string) bool {
				if _, err := os.Stat(filepath.Join(bundleDir, "examplebun-0.1.0")); os.IsNotExist(err) {
					return false
				}
				return true
			},
		},
		"includes --destination": {
			dest: filepath.Join(tempDir, "unzipped"), // Example directory for unzipped bundles
			fileCreated: func(bundleDir string) bool {
				if _, err := os.Stat(filepath.Join(bundleDir, "examplebun-0.1.0")); os.IsNotExist(err) {
					return false
				}
				return true
			},
		},
	}

	for name, testCase := range tests {
		func() {
			testHome := CreateTestHome(t)
			defer os.RemoveAll(testHome.String())

			impCmd := &importCmd{
				source:  "testdata/import/examplebun-0.1.0.tgz",
				dest:    testCase.dest,
				out:     ioutil.Discard,
				home:    testHome,
				verbose: false,
			}

			err = impCmd.run()
			if err != nil {
				t.Fatalf("Fail on test: %s, error message: %s", name, err)
			}
			is.Equal(testCase.dest != "", testCase.fileCreated(impCmd.dest), "Fail on test: "+name)
		}()
	}

}
