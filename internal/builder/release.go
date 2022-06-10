// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
)

type Release struct {
	BuildOutputDirectory   string
	ReleaseOutputDirectory string
}

func (r *Release) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Releasing chaincode...")

	// TODO is this required?
	// imageSrcPath := filepath.Join(r.BuildOutputDirectory, "image.json")
	// imageDestPath := filepath.Join(r.ReleaseOutputDirectory, "image.json")
	// err := copy.Copy(imageSrcPath, imageDestPath)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error copying %s to %s: %s\n", imageSrcPath, imageDestPath, err)
	// 	return err
	// }

	// TODO copy any META-INF
	// metainfSrcPath := filepath.Join(r.BuildOutputDirectory, "META-INF")
	// metainfDestPath := filepath.Join(r.ReleaseOutputDirectory, "META-INF")

	return nil
}
