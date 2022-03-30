// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Detect struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
}

type metadata struct {
	Type string `json:"type"`
}

func (d *Detect) Run() error {
	mdpath := filepath.Join(d.ChaincodeMetadataDirectory, "metadata.json")
	mdbytes, err := ioutil.ReadFile(mdpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s", mdpath, err)
		return err
	}

	var metadata metadata
	err = json.Unmarshal(mdbytes, &metadata)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading json %s: %s", mdpath, err)
		return err
	}

	if strings.ToLower(metadata.Type) == "k8s" {
		return nil
	}

	return fmt.Errorf("chaincode type not supported: %s", metadata.Type)
}
