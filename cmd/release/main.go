// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
)

const (
	expectedArgsLength        = 3
	buildOutputDirectoryArg   = 1
	releaseOutputDirectoryArg = 2
)

func main() {
	debug := util.GetOptionalEnv(util.DebugVariable, "false")
	ctx := log.NewCmdContext(context.Background(), debug == "true")
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println("Expected BUILD_OUTPUT_DIR and RELEASE_OUTPUT_DIR arguments")
		os.Exit(1)
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
		os.Exit(1)
	}
}
