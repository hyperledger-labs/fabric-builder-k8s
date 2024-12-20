// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

func Build() int {
	const (
		expectedArgsLength            = 4
		chaincodeSourceDirectoryArg   = 1
		chaincodeMetadataDirectoryArg = 2
		buildOutputDirectoryArg       = 3
	)

	debug, _ := strconv.ParseBool(util.GetOptionalEnv(util.DebugVariable, "false"))
	ctx := log.NewCmdContext(context.Background(), debug)
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println(
			"Expected CHAINCODE_SOURCE_DIR, CHAINCODE_METADATA_DIR and BUILD_OUTPUT_DIR arguments",
		)

		return 1
	}

	chaincodeSourceDirectory := os.Args[chaincodeSourceDirectoryArg]
	chaincodeMetadataDirectory := os.Args[chaincodeMetadataDirectoryArg]
	buildOutputDirectory := os.Args[buildOutputDirectoryArg]

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

		return 1
	}

	return 0
}
