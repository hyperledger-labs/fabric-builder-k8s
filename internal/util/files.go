// SPDX-License-Identifier: Apache-2.0

package util

import v1 "k8s.io/api/core/v1"

// ChaincodeJson represents the chaincode.json file that is supplied by Fabric in
// the RUN_METADATA_DIR
type ChaincodeJson struct {
	ChaincodeID string `json:"chaincode_id"`
	PeerAddress string `json:"peer_address"`
	ClientCert  string `json:"client_cert"`
	ClientKey   string `json:"client_key"`
	RootCert    string `json:"root_cert"`
	MspID       string `json:"mspid"`
}

// image represents the image.json file from the k8s chaincode tarball
type ChaincodeImage struct {
	Name            string        `json:"name"`
	Tag             string        `json:"tag"`
	ImagePullPolicy v1.PullPolicy `json:"imagePullPolicy,omitempty"`
}
