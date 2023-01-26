// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
)

type Detect struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
}

type metadata struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

var ErrUnsupportedChaincodeType = errors.New("chaincode type not supported")

func (d *Detect) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Checking chaincode type...")

	mdpath := filepath.Join(d.ChaincodeMetadataDirectory, "metadata.json")

	mdbytes, err := os.ReadFile(mdpath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", mdpath, err)
	}

	var metadata metadata

	err = json.Unmarshal(mdbytes, &metadata)
	if err != nil {
		return fmt.Errorf("unable to process %s: %w", mdpath, err)
	}

	if strings.ToLower(metadata.Type) == "k8s" {
		logger.Printf("Detected k8s chaincode: %s", metadata.Label)

		return nil
	}

	logger.Debugf("Chaincode type not supported: %s", metadata.Type)

	return ErrUnsupportedChaincodeType
}
