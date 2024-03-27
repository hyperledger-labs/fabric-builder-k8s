// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/otiai10/copy"
)

// CopyImageJSON validates and copies the chaincode image file.
func CopyImageJSON(logger *log.CmdLogger, src, dest string) error {
	imageSrcPath := filepath.Join(src, ImageFile)
	imageDestPath := filepath.Join(dest, ImageFile)

	logger.Debugf("Copying chaincode image file from %s to %s", imageSrcPath, imageDestPath)

	err := copy.Copy(imageSrcPath, imageDestPath)
	if err != nil {
		return fmt.Errorf(
			"failed to copy chaincode image file from %s to %s: %w",
			imageSrcPath,
			imageDestPath,
			err,
		)
	}

	logger.Debugf("Verifying chaincode image file %s", imageDestPath)

	_, err = ReadImageJSON(logger, dest)
	if err != nil {
		return err
	}

	return nil
}

// CopyIndexFiles copies CouchDB index definitions from source to destination directories.
func CopyIndexFiles(logger *log.CmdLogger, src, dest string) error {
	indexDir := filepath.Join("statedb", "couchdb")
	indexSrcDir := filepath.Join(src, MetadataDir, indexDir)
	indexDestDir := filepath.Join(dest, indexDir)

	logger.Debugf("Copying CouchDB index definitions from %s to %s", indexSrcDir, indexDestDir)

	opt := copy.Options{
		Skip: func(info os.FileInfo, src, _ string) (bool, error) {
			logger.Debugf("Checking source copy path: %s", src)

			if info.IsDir() {
				skip, err := skipFolder(logger, indexSrcDir, src)

				if err != nil {
					return skip, fmt.Errorf(
						"error checking if the folder is eligible to have a couchdb index: %s, %s: %w",
						indexSrcDir,
						src,
						err,
					)
				}

				return skip, nil
			}

			logger.Debugf("Checking if it is a JSON file: %s", src)

			return !strings.HasSuffix(src, ".json"), nil
		},
	}

	if err := copy.Copy(indexSrcDir, indexDestDir, opt); err != nil {
		return fmt.Errorf(
			"failed to copy CouchDB index definitions from %s to %s: %w",
			indexSrcDir,
			indexDestDir,
			err,
		)
	}

	return nil
}

// skipFolder checks if the folder will need to be skipped during indexes copy.
func skipFolder(logger *log.CmdLogger, indexSrcDir, src string) (bool, error) {
	path, err := filepath.Rel(indexSrcDir, src)

	if err != nil {
		return true, fmt.Errorf(
			"error resolving the relative path: %s, %s: %w",
			indexSrcDir,
			src,
			err,
		)
	}

	matchContainsPublicIndexFolder, err := filepath.Match("indexes", path)

	if err != nil {
		return true, fmt.Errorf(
			"error matching the path with public index couchdb folder: %s: %w",
			path,
			err,
		)
	}

	matchContainsPrivateDataCollectionFolder, err := filepath.Match("collections", path)

	if err != nil {
		return true, fmt.Errorf(
			"error matching the path with the collection folder: %s: %w",
			path,
			err,
		)
	}

	matchPrivateDataCollectionFolder, err := filepath.Match("collections/*", path)

	if err != nil {
		return true, fmt.Errorf(
			"error matching the path with the private data collection definition folder: %s: %w",
			path,
			err,
		)
	}

	matchPrivateDataCollectionIndexFolder, err := filepath.Match("collections/*/indexes", path)

	if err != nil {
		return true, fmt.Errorf(
			"error matching the path with the private data collection index definition folder: %s: %w",
			path,
			err,
		)
	}
	relativeFoldersLength := len(strings.Split(path, string(filepath.Separator)))

	logger.Debugf("relative path: %s, total relative folders: %d", path, relativeFoldersLength)
	logger.Debugf("Match pattern - index: %t", matchContainsPublicIndexFolder)
	logger.Debugf("Match pattern - collections: %t", matchContainsPrivateDataCollectionFolder)
	logger.Debugf("Match pattern - collections/*: %t", matchPrivateDataCollectionFolder)
	logger.Debugf("Match pattern - collections/*/indexes: %t", matchPrivateDataCollectionIndexFolder)

	switch {
	case relativeFoldersLength == 1 && (!matchContainsPublicIndexFolder && !matchContainsPrivateDataCollectionFolder):
		logger.Debugf("Should skip folder")

		return true, nil

	case relativeFoldersLength == 2 && (!matchPrivateDataCollectionFolder):
		logger.Debugf("Should skip folder")

		return true, nil

	case relativeFoldersLength == 3 && (!matchPrivateDataCollectionIndexFolder):
		logger.Debugf("Should skip folder")

		return true, nil

	default:
		logger.Debugf("Should NOT skip folder")

		return false, nil
	}
}

// CopyMetadataDir copies all chaincode metadata from source to destination directories.
func CopyMetadataDir(logger *log.CmdLogger, src, dest string) error {
	metadataSrcDir := filepath.Join(src, MetadataDir)
	metadataDestDir := filepath.Join(dest, MetadataDir)

	logger.Debugf("Copying chaincode metadata from %s to %s", metadataSrcDir, metadataDestDir)

	fileInfo, err := os.Lstat(metadataSrcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// metadata is optional
			return nil
		}

		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("chaincode metadata path %s is not a directory: %w", metadataSrcDir, err)
	}

	if err := copy.Copy(metadataSrcDir, metadataDestDir); err != nil {
		return fmt.Errorf(
			"failed to copy chaincode metadata from %s to %s: %w",
			metadataSrcDir,
			metadataDestDir,
			err,
		)
	}

	return nil
}
