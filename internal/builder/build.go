// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

type Build struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
	BuildOutputDirectory       string
}

func (b *Build) Run() error {
	imageSrcPath := filepath.Join(b.ChaincodeSourceDirectory, "image.json")
	imageDestPath := filepath.Join(b.BuildOutputDirectory, "image.json")
	err := copy.Copy(imageSrcPath, imageDestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying %s to %s: %s\n", imageSrcPath, imageDestPath, err)
		return err
	}

	// TODO copy any META-INF
	// metainfSrcPath := filepath.Clean(filepath.Join(sourceDir, "META-INF"))
	// metainfDestPath := filepath.Clean(filepath.Join(outputDir, "META-INF"))

	return nil
}
