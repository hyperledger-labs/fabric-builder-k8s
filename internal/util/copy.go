// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/otiai10/copy"
)

// CopyImageJson validates and copies the chaincode image file
func CopyImageJson(logger *log.CmdLogger, src, dest string) error {
	imageSrcPath := filepath.Join(src, ImageFile)
	imageDestPath := filepath.Join(dest, ImageFile)

	logger.Debugf("Copying chaincode image file from %s to %s", imageSrcPath, imageDestPath)
	err := copy.Copy(imageSrcPath, imageDestPath)
	if err != nil {
		return fmt.Errorf("failed to copy chaincode image file from %s to %s: %w", imageSrcPath, imageDestPath, err)
	}

	logger.Debugf("Verifying chaincode image file %s", imageDestPath)
	_, err = ReadImageJson(logger, dest)
	if err != nil {
		return err
	}

	return nil
}

// CopyIndexFiles copies CouchDB index definitions from source to destination directories
func CopyIndexFiles(logger *log.CmdLogger, src, dest string) error {
	indexDir := filepath.Join("statedb", "couchdb", "indexes")
	indexSrcDir := filepath.Join(src, MetadataDir, indexDir)
	indexDestDir := filepath.Join(dest, indexDir)

	logger.Debugf("Copying CouchDB index definitions from %s to %s", indexSrcDir, indexDestDir)
	fi, err := os.Lstat(indexSrcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// indexes are optional
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("CouchDB index definitions path %s is not a directory: %w", indexSrcDir, err)
	}

	opt := copy.Options{
		Skip: func(src string) (bool, error) {
			return !strings.HasSuffix(src, ".json"), nil
		},
	}
	if err := copy.Copy(indexSrcDir, indexDestDir, opt); err != nil {
		return fmt.Errorf("failed to copy CouchDB index definitions from %s to %s: %w", indexSrcDir, indexDestDir, err)
	}

	return nil
}

// CopyMetadataDir copies all chaincode metadata from source to destination directories
func CopyMetadataDir(logger *log.CmdLogger, src, dest string) error {
	metadataSrcDir := filepath.Join(src, MetadataDir)
	metadataDestDir := filepath.Join(dest, MetadataDir)

	logger.Debugf("Copying chaincode metadata from %s to %s", metadataSrcDir, metadataDestDir)
	fi, err := os.Lstat(metadataSrcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// metadata is optional
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("chaincode metadata path %s is not a directory: %w", metadataSrcDir, err)
	}

	if err := copy.Copy(metadataSrcDir, metadataDestDir); err != nil {
		return fmt.Errorf("failed to copy chaincode metadata from %s to %s: %w", metadataSrcDir, metadataDestDir, err)
	}

	return nil
}
