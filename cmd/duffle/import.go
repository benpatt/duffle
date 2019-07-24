package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/deislabs/duffle/pkg/duffle/home"
	"github.com/deislabs/duffle/pkg/loader"
	"github.com/deislabs/duffle/pkg/packager"
)

const importDesc = `
Imports a Cloud Native Application Bundle (CNAB) from a gzipped tar file on the local file system.

This operation is symmetric with export, and so will save the bundle, along with any images 
referenced by the bundle, to the internal storage.

If --destination is specified, the file will also be unzipped to the specified path. Please
note -d with the empty string "" will be ignored, so use '-d .' to save the unzipped bundle
to the working directory.

Example:
	$ duffle import mybundle.tgz
	$ duffle import mybundle.tgz -d bundles/unzipped
`

type importCmd struct {
	source  string
	dest    string
	out     io.Writer
	home    home.Home
	verbose bool
}

func newImportCmd(w io.Writer) *cobra.Command {
	importc := &importCmd{
		out:  w,
		home: home.Home(homePath()),
	}

	cmd := &cobra.Command{
		Use:   "import [PATH]",
		Short: "save CNAB bundle from gzipped tar file",
		Long:  importDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("this command requires the path to the packaged bundle")
			}
			importc.source = args[0]

			return importc.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&importc.dest, "destination", "d", "", "Location to unpack bundle")
	f.BoolVarP(&importc.verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func (im *importCmd) run() error {
	source, err := filepath.Abs(im.source)
	if err != nil {
		return fmt.Errorf("Error in source: %s", err)
	}

	var dest string
	if im.dest == "" {
		dest, err = ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dest)
	} else {
		dest, err = filepath.Abs(im.dest)
		if err != nil {
			return fmt.Errorf("Error in destination: %s", err)
		}
	}

	l := loader.NewLoader()
	imp, err := packager.NewImporter(source, dest, l, im.verbose)
	if err != nil {
		return err
	}
	return imp.Import(im.home)
}
