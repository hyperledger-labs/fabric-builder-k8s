// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	"k8s.io/apimachinery/pkg/util/validation"
)

type Build struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
	BuildOutputDirectory       string
}

func (b *Build) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Building chaincode...")

	metadata, err := util.ReadMetadataJSON(logger, b.ChaincodeMetadataDirectory)
	if err != nil {
		return err
	}

	if errs := validation.IsDNS1035Label(metadata.Label); len(errs) != 0 {
		return fmt.Errorf(
			"chaincode label '%s' must be a valid RFC1035 label: %v",
			metadata.Label,
			errs,
		)
	}

	err = util.CopyImageJSON(logger, b.ChaincodeSourceDirectory, b.BuildOutputDirectory)
	if err != nil {
		return err
	}

	err = util.CopyMetadataDir(logger, b.ChaincodeSourceDirectory, b.BuildOutputDirectory)
	if err != nil {
		return err
	}

	return nil
}
