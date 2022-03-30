// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "Expected CHAINCODE_SOURCE_DIR, CHAINCODE_METADATA_DIR and BUILD_OUTPUT_DIR arguments")
		os.Exit(1)
	}

	build := &builder.Build{
		ChaincodeSourceDirectory:   os.Args[1],
		ChaincodeMetadataDirectory: os.Args[2],
		BuildOutputDirectory:       os.Args[3],
	}

	if err := build.Run(); err != nil {
		os.Exit(1)
	}
}
