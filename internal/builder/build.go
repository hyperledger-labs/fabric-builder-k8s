// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
)

type Build struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
	BuildOutputDirectory       string
}

func (b *Build) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Building chaincode...")

	err := util.CopyImageJson(logger, b.ChaincodeSourceDirectory, b.BuildOutputDirectory)
	if err != nil {
		return err
	}

	err = util.CopyMetadataDir(logger, b.ChaincodeSourceDirectory, b.BuildOutputDirectory)
	if err != nil {
		return err
	}

	return nil
}
