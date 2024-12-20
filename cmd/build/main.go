// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/hyperledger-labs/fabric-builder-k8s/cmd"
)

func main() {
	os.Exit(cmd.Build())
}
