// SPDX-License-Identifier: Apache-2.0

package builder

import "errors"

type Run struct {
	BuildOutputDirectory string
	RunMetadataDirectory string
}

func (d *Run) Run() error {
	return errors.New("not implemented yet")
}
