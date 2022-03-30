// SPDX-License-Identifier: Apache-2.0

package builder

import "errors"

type Release struct {
	BuildOutputDirectory   string
	ReleaseOutputDirectory string
}

func (d *Release) Run() error {
	return errors.New("not implemented yet")
}
