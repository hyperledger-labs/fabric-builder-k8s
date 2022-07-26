// SPDX-License-Identifier: Apache-2.0

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
)

// ChaincodeJSON represents the chaincode.json file that is supplied by Fabric in
// the RUN_METADATA_DIR.
type ChaincodeJSON struct {
	ChaincodeID string `json:"chaincode_id"`
	PeerAddress string `json:"peer_address"`
	ClientCert  string `json:"client_cert"`
	ClientKey   string `json:"client_key"`
	RootCert    string `json:"root_cert"`
	MspID       string `json:"mspid"`
}

// ImageJSON represents the image.json file in the k8s chaincode package.
type ImageJSON struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
}

const (
	ChaincodeFile = "chaincode.json"
	ImageFile     = "image.json"
	MetadataDir   = "META-INF"
)

// ReadChaincodeJSON reads and parses the chaincode.json file in the provided directory.
func ReadChaincodeJSON(logger *log.CmdLogger, dir string) (*ChaincodeJSON, error) {
	chaincodeJSONPath := filepath.Join(dir, ChaincodeFile)
	logger.Debugf("Reading %s...", chaincodeJSONPath)

	chaincodeJSONContents, err := os.ReadFile(chaincodeJSONPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", chaincodeJSONPath, err)
	}

	var chaincodeData ChaincodeJSON
	if err := json.Unmarshal(chaincodeJSONContents, &chaincodeData); err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", chaincodeJSONPath, err)
	}

	logger.Debugf("Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	return &chaincodeData, nil
}

// ReadImageJSON reads and parses the image.json file in the provided directory.
func ReadImageJSON(logger *log.CmdLogger, dir string) (*ImageJSON, error) {
	imageJSONPath := filepath.Join(dir, ImageFile)
	logger.Debugf("Reading %s...", imageJSONPath)

	imageJSONContents, err := os.ReadFile(imageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", imageJSONPath, err)
	}

	var imageData ImageJSON
	if err := json.Unmarshal(imageJSONContents, &imageData); err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", imageJSONPath, err)
	}

	logger.Debugf("Image name: %s\nImage digest: %s\n", imageData.Name, imageData.Digest)

	if len(imageData.Name) == 0 || len(imageData.Digest) == 0 {
		return nil, fmt.Errorf("%s file must contain 'name' and 'digest'", imageJSONPath)
	}

	return &imageData, nil
}
