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

func Release() int {
	const (
		expectedArgsLength        = 3
		buildOutputDirectoryArg   = 1
		releaseOutputDirectoryArg = 2
	)

	debug, _ := strconv.ParseBool(util.GetOptionalEnv(util.DebugVariable, "false"))
	ctx := log.NewCmdContext(context.Background(), debug)
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println("Expected BUILD_OUTPUT_DIR and RELEASE_OUTPUT_DIR arguments")

		return 1
	}

	buildOutputDirectory := os.Args[buildOutputDirectoryArg]
	releaseOutputDirectory := os.Args[releaseOutputDirectoryArg]

	logger.Debugf("Build output directory: %s", buildOutputDirectory)
	logger.Debugf("Release output directory: %s", releaseOutputDirectory)

	release := &builder.Release{
		BuildOutputDirectory:   buildOutputDirectory,
		ReleaseOutputDirectory: releaseOutputDirectory,
	}

	if err := release.Run(ctx); err != nil {
		logger.Printf("Error releasing chaincode: %+v", err)

		return 1
	}

	return 0
}
