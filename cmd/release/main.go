// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Expected BUILD_OUTPUT_DIR and RELEASE_OUTPUT_DIR arguments")
		os.Exit(1)
	}

	release := &builder.Release{
		BuildOutputDirectory:   os.Args[1],
		ReleaseOutputDirectory: os.Args[2],
	}

	if err := release.Run(); err != nil {
		os.Exit(1)
	}
}
