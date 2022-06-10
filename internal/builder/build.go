// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/otiai10/copy"
)

type Build struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
	BuildOutputDirectory       string
}

func (b *Build) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Building chaincode...")

	imageSrcPath := filepath.Join(b.ChaincodeSourceDirectory, "image.json")
	imageDestPath := filepath.Join(b.BuildOutputDirectory, "image.json")
	err := copy.Copy(imageSrcPath, imageDestPath)
	if err != nil {
		return fmt.Errorf("could not copy %s to %s: %w", imageSrcPath, imageDestPath, err)
	}

	// TODO copy any META-INF
	// metainfSrcPath := filepath.Join(b.ChaincodeSourceDirectory, "META-INF")
	// metainfDestPath := filepath.Join(b.BuildOutputDirectory, "META-INF")

	return nil
}
