package packager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/deislabs/duffle/pkg/duffle/home"

	"github.com/stretchr/testify/assert"

	"github.com/deislabs/duffle/pkg/loader"
)

// TODO: Make a central "create temp home" for use in command tests across the project.
//  See main_test.go, export.go, and there are probably more occurrences in other tests
func createTempHome() (home.Home, error) {
	tempDir, err := ioutil.TempDir("", "temp-home")
	tempHome := home.Home(tempDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(tempHome.Bundles(), 0755); err != nil {
		defer os.RemoveAll(tempHome.String())
		return "", err
	}
	f, err := os.OpenFile(tempHome.Repositories(), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		defer os.RemoveAll(tempHome.String())
		return "", err
	}
	f.Close()
	return tempHome, nil
}

func TestImport(t *testing.T) {
	is := assert.New(t)

	tests := map[string]struct {
		dest         string
		expectsError bool
		validate     func(string, home.Home) bool
	}{
		"Well formed bundle": {
			dest:         "testdata/examplebun-0.1.0.tgz",
			expectsError: false,
			validate: func(tempDir string, home home.Home) bool {
				// Check for unzipped file
				if _, err := os.Stat(filepath.Join(tempDir, "examplebun-0.1.0")); os.IsNotExist(err) {
					return false
				}
				// Check file has been added to internal storage bundles
				bundles, err := ioutil.ReadDir(home.Bundles())
				if err != nil {
					t.Fatal(err)
				}
				if len(bundles) != 1 {
					return false
				}
				// Check file has been added to repos json
				repos, err := ioutil.ReadFile(home.Repositories())
				if err != nil {
					t.Fatal(err)
				}
				if string(repos) == "" {
					return false
				}
				return true
			},
		},
		"Malformed bundle": {
			dest:         "testdata/malformed-0.1.0.tgz",
			expectsError: true,
			validate: func(tempDir string, home home.Home) bool {
				// Check for unzipped file, should not exist
				if _, err := os.Stat(filepath.Join(tempDir, "malformed-0.1.0")); !os.IsNotExist(err) {
					return false
				}
				// Check file has not been added to internal storage bundles
				bundles, err := ioutil.ReadDir(home.Bundles())
				if err != nil {
					t.Fatal(err)
				}
				if len(bundles) != 0 {
					return false
				}
				// Check file has not been added to repos json
				repos, err := ioutil.ReadFile(home.Repositories())
				if err != nil {
					t.Fatal(err)
				}
				if string(repos) != "" {
					return false
				}
				return true
			},
		},
	}

	for name, testCase := range tests {
		func() {
			tempHome, err := createTempHome()
			if err != nil {
				t.Fatal(fmt.Errorf("problem creating temp tempHome: %s", err))
			}
			defer os.RemoveAll(tempHome.String())

			tempDir, err := ioutil.TempDir("", "duffle-import-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			im := &Importer{
				Source:      testCase.dest,
				Destination: tempDir,
				Loader:      loader.NewLoader(),
			}

			if testCase.expectsError {
				is.Error(im.Import(tempHome), "No error on import of test: "+name)
			} else {
				is.NoError(im.Import(tempHome), "Error on import of test: "+name)
			}
			is.True(testCase.validate(tempDir, tempHome), "Validation fail on test: "+name)
		}()
	}
}
