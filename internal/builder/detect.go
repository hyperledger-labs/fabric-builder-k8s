// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
)

type Detect struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
}

type metadata struct {
	Type string `json:"type"`
}

func (d *Detect) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Checking chaincode type...")

	mdpath := filepath.Join(d.ChaincodeMetadataDirectory, "metadata.json")
	mdbytes, err := ioutil.ReadFile(mdpath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", mdpath, err)
	}

	var metadata metadata
	err = json.Unmarshal(mdbytes, &metadata)
	if err != nil {
		return fmt.Errorf("unable to process %s: %w", mdpath, err)
	}

	if strings.ToLower(metadata.Type) == "k8s" {
		return nil
	}

	return fmt.Errorf("chaincode type not supported: %s", metadata.Type)
}
