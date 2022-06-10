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

	if len(os.Args) != 3 {
		logger.Println("Expected CHAINCODE_SOURCE_DIR and CHAINCODE_METADATA_DIR arguments")
		os.Exit(1)
	}
	chaincodeSourceDirectory := os.Args[1]
	chaincodeMetadataDirectory := os.Args[2]
	logger.Debugf("Chaincode source directory: %s", chaincodeSourceDirectory)
	logger.Debugf("Chaincode metadata directory: %s", chaincodeMetadataDirectory)

	detect := &builder.Detect{
		ChaincodeSourceDirectory:   chaincodeSourceDirectory,
		ChaincodeMetadataDirectory: chaincodeMetadataDirectory,
	}

	if err := detect.Run(ctx); err != nil {
		logger.Printf("Error detecting chaincode: %+v", err)
		os.Exit(1)
	}
}
