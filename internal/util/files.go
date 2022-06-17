// SPDX-License-Identifier: Apache-2.0

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
)

// ChaincodeJson represents the chaincode.json file that is supplied by Fabric in
// the RUN_METADATA_DIR
type ChaincodeJson struct {
	ChaincodeID string `json:"chaincode_id"`
	PeerAddress string `json:"peer_address"`
	ClientCert  string `json:"client_cert"`
	ClientKey   string `json:"client_key"`
	RootCert    string `json:"root_cert"`
	MspID       string `json:"mspid"`
}

// ImageJson represents the image.json file in the k8s chaincode package
type ImageJson struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
}

const (
	ChaincodeFile = "chaincode.json"
	ImageFile     = "image.json"
	MetadataDir   = "META-INF"
)

// ReadChaincodeJson reads and parses the chaincode.json file in the provided directory
func ReadChaincodeJson(logger *log.CmdLogger, dir string) (*ChaincodeJson, error) {
	chaincodeJsonPath := filepath.Join(dir, ChaincodeFile)
	logger.Debugf("Reading %s...", chaincodeJsonPath)

	chaincodeJsonContents, err := os.ReadFile(chaincodeJsonPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", chaincodeJsonPath, err)
	}

	var chaincodeData ChaincodeJson
	if err := json.Unmarshal(chaincodeJsonContents, &chaincodeData); err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", chaincodeJsonPath, err)
	}

	logger.Debugf("Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	return &chaincodeData, nil
}

// ReadImageJson reads and parses the image.json file in the provided directory
func ReadImageJson(logger *log.CmdLogger, dir string) (*ImageJson, error) {
	imageJsonPath := filepath.Join(dir, ImageFile)
	logger.Debugf("Reading %s...", imageJsonPath)

	imageJsonContents, err := os.ReadFile(imageJsonPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", imageJsonPath, err)
	}

	var imageData ImageJson
	if err := json.Unmarshal(imageJsonContents, &imageData); err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", imageJsonPath, err)
	}

	logger.Debugf("Image name: %s\nImage digest: %s\n", imageData.Name, imageData.Digest)

	if len(imageData.Name) == 0 || len(imageData.Digest) == 0 {
		return nil, fmt.Errorf("%s file must contain 'name' and 'digest'", imageJsonPath)
	}

	return &imageData, nil
}
