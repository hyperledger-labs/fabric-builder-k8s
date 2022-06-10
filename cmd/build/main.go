// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
)

func main() {
	debug := util.GetOptionalEnv(util.DebugVariable, "false")
	ctx := log.NewCmdContext(context.Background(), debug == "true")
	logger := log.New(ctx)

	if len(os.Args) != 4 {
		logger.Println("Expected CHAINCODE_SOURCE_DIR, CHAINCODE_METADATA_DIR and BUILD_OUTPUT_DIR arguments")
		os.Exit(1)
	}
	chaincodeSourceDirectory := os.Args[1]
	chaincodeMetadataDirectory := os.Args[2]
	buildOutputDirectory := os.Args[3]
	logger.Debugf("Chaincode source directory: %s", chaincodeSourceDirectory)
	logger.Debugf("Chaincode metadata directory: %s", chaincodeMetadataDirectory)
	logger.Debugf("Build output directory: %s", buildOutputDirectory)

	build := &builder.Build{
		ChaincodeSourceDirectory:   chaincodeSourceDirectory,
		ChaincodeMetadataDirectory: chaincodeMetadataDirectory,
		BuildOutputDirectory:       buildOutputDirectory,
	}

	if err := build.Run(ctx); err != nil {
		logger.Printf("Error building chaincode: %+v", err)
		os.Exit(1)
	}
}
