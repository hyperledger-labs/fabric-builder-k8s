// SPDX-License-Identifier: Apache-2.0

package util

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

// ImageJson represents the image.json file in the k8s chaincode package
type ImageJson struct {
	Name   string `json:"name"`
	Digest string `json:"digest"`
}
