// SPDX-License-Identifier: Apache-2.0

package util

import "strings"

type ChaincodePackageID struct {
	Label string
	Hash  string
}

// NewChaincodePackageID returns a ChaincodePackageID created from the provided string.
func NewChaincodePackageID(chaincodeID string) *ChaincodePackageID {
	substrings := strings.Split(chaincodeID, ":")

	// If it doesn't look like a label and a hash, don't try and guess which is which
	if len(substrings) == 1 {
		return &ChaincodePackageID{
			Label: "",
			Hash:  "",
		}
	}

	return &ChaincodePackageID{
		Label: strings.Join(substrings[:len(substrings)-1], ":"),
		Hash:  substrings[len(substrings)-1],
	}
}
