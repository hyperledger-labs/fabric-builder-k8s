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

	logger.Debugf("Copying couchdb index files from %s to %s", indexSrcDir, indexDestDir)

	_, err := os.Lstat(indexSrcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// indexes are optional
			return nil
		}

		return err
	}

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

			skip, err := skipFile(logger, indexSrcDir, src)
			if err != nil {
				return skip, fmt.Errorf(
					"error checking if the file is eligible to have a couchdb index: %s, %s: %w",
					indexSrcDir,
					src,
					err,
				)
			}

			return skip, nil
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

// skipFile checks if the file will need to be skipped during indexes copy.
func skipFile(logger *log.CmdLogger, indexSrcDir, src string) (bool, error) {
	path, err := filepath.Rel(indexSrcDir, src)
	if err != nil {
		logger.Debugf("error verifying relative path from: %s, src: %s", indexSrcDir, src)

		return true, fmt.Errorf(
			"error verifying relative path from %s to %s: %w",
			indexSrcDir,
			src,
			err,
		)
	}

	if len(strings.Split(path, string(filepath.Separator))) == 1 { // JSON is in root couchdb folder
		logger.Debugf("The JSON file in the root couchdb index folder, should skip: %s, src: %s", path, src)

		return true, nil
	}

	if strings.HasSuffix(src, ".json") {
		logger.Debugf("The JSON file is valid, should copy: %s, src: %s", path, src)

		return false, nil
	}

	logger.Debugf("The JSON file is invalid, should skip: %s, src: %s", path, src)

	return true, nil
}

// skipFolder checks if the folder will need to be skipped during indexes copy.
func skipFolder(logger *log.CmdLogger, indexSrcDir, src string) (bool, error) {
	path, err := filepath.Rel(indexSrcDir, src)
	if err != nil {
		logger.Debugf("failed resolve relative path: %s, src: %s", indexSrcDir, src)

		return true, fmt.Errorf("failed resolve relative path %s to %s: %w", indexSrcDir, src, err)
	}

	matchContainsPublicIndexFolder, _ := filepath.Match("indexes", path)
	matchContainsPrivateDataCollectionFolder, _ := filepath.Match("collections", path)
	matchPrivateDataCollectionFolder, _ := filepath.Match("collections/*", path)
	matchPrivateDataCollectionIndexFolder, _ := filepath.Match("collections/*/indexes", path)
	relativeFoldersLength := len(strings.Split(path, string(filepath.Separator)))

	logger.Debugf("Calculated relative path: %s. Total relative folders: %d", path, relativeFoldersLength)
	logger.Debugf("Match pattern 'index': %t", matchContainsPublicIndexFolder)
	logger.Debugf("Match pattern 'collections': %t", matchContainsPrivateDataCollectionFolder)
	logger.Debugf("Match pattern 'collections/*': %t", matchPrivateDataCollectionFolder)
	logger.Debugf("Match pattern 'collections/*/indexes': %t", matchPrivateDataCollectionIndexFolder)

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
		logger.Debugf("Should not skip folder")

		return false, nil
	}
}
