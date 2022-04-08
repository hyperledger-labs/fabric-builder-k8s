// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Expected BUILD_OUTPUT_DIR and RUN_METADATA_DIR arguments")
		os.Exit(1)
	}

	run := &builder.Run{
		BuildOutputDirectory: os.Args[1],
		RunMetadataDirectory: os.Args[2],
	}

	if err := run.Run(); err != nil {
		// TODO better error handling?
		fmt.Println(err)
		os.Exit(1)
	}
}
