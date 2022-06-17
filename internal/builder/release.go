// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
)

type Release struct {
	BuildOutputDirectory   string
	ReleaseOutputDirectory string
}

func (r *Release) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Releasing chaincode...")

	// If CouchDB index definitions are required for the chaincode, release is
	// responsible for placing the indexes into the statedb/couchdb/indexes
	// directory under RELEASE_OUTPUT_DIR. The indexes must have a .json
	// extension.
	err := util.CopyIndexFiles(logger, r.BuildOutputDirectory, r.ReleaseOutputDirectory)
	if err != nil {
		return err
	}

	return nil
}
