// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Expected CHAINCODE_SOURCE_DIR and CHAINCODE_METADATA_DIR arguments")
		os.Exit(1)
	}

	detect := &builder.Detect{
		ChaincodeSourceDirectory:   os.Args[1],
		ChaincodeMetadataDirectory: os.Args[2],
	}

	if err := detect.Run(); err != nil {
		os.Exit(1)
	}
}
