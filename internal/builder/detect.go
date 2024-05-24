// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"errors"
	"strings"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

type Detect struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
}

var ErrUnsupportedChaincodeType = errors.New("chaincode type not supported")

func (d *Detect) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Checking chaincode type...")

	metadata, err := util.ReadMetadataJSON(logger, d.ChaincodeMetadataDirectory)
	if err != nil {
		return err
	}

	if strings.ToLower(metadata.Type) == "k8s" {
		logger.Printf("Detected k8s chaincode: %s", metadata.Label)

		return nil
	}

	logger.Debugf("Chaincode type not supported: %s", metadata.Type)

	return ErrUnsupportedChaincodeType
}
