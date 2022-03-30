// SPDX-License-Identifier: Apache-2.0

package builder

import "errors"

type Build struct {
	ChaincodeSourceDirectory   string
	ChaincodeMetadataDirectory string
	BuildOutputDirectory       string
}

func (d *Build) Run() error {
	return errors.New("not implemented yet")
}
