// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/builder"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Expected BUILD_OUTPUT_DIR and RUN_METADATA_DIR arguments")
		os.Exit(1)
	}

	peerID, err := util.GetRequiredEnv("CORE_PEER_ID")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Expected CORE_PEER_ID environment variable")
		os.Exit(1)
	}

	kubeconfigPath := util.GetOptionalEnv("KUBECONFIG_PATH", "")

	kubeNamespace := util.GetOptionalEnv("KUBE_NAMESPACE", "")
	if kubeNamespace == "" {
		kubeNamespace, err = util.GetKubeNamespace()
		if err != nil {
			kubeNamespace = "default"
		}
	}

	run := &builder.Run{
		BuildOutputDirectory: os.Args[1],
		RunMetadataDirectory: os.Args[2],
		PeerID:               peerID,
		KubeconfigPath:       kubeconfigPath,
		KubeNamespace:        kubeNamespace,
	}

	if err := run.Run(); err != nil {
		// TODO better error handling?
		fmt.Fprintf(os.Stderr, "Error running chaincode.\nBuild dir: %s\nRun dir: %s\nError: %v\n", run.BuildOutputDirectory, run.RunMetadataDirectory, err)
		os.Exit(1)
	}

	// TODO nasty hack to keep chaincode running- peer assumes chaincode has terminated when builder run terminates!
	var m sync.Mutex
	m.Lock()
	m.Lock()
}
