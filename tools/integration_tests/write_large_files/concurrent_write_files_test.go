// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package write_large_files

import (
	"fmt"
	"os"
	"path"
	"syscall"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/v2/tools/integration_tests/util/operations"
	"github.com/googlecloudplatform/gcsfuse/v2/tools/integration_tests/util/setup"
	"golang.org/x/sync/errgroup"
)

const (
	DirForConcurrentWrite = "dirForConcurrentWrite"
)

var FileOne = "fileOne" + setup.GenerateRandomString(5) + ".txt"
var FileTwo = "fileTwo" + setup.GenerateRandomString(5) + ".txt"
var FileThree = "fileThree" + setup.GenerateRandomString(5) + ".txt"

func writeFile(filePath string, fileSize int64) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|syscall.O_DIRECT, WritePermission_0200)
	if err != nil {
		return fmt.Errorf("Open file for write at start: %v", err)
	}

	// Closing file at the end.
	defer operations.CloseFile(f)
	return operations.WriteChunkOfRandomBytesToFile(f, int(fileSize), 0)
}

func validateFileContents(fileName string, mntFilePath string, t *testing.T) error {
	filePathInGcsBucket := path.Join(DirForConcurrentWrite, fileName)
	localFilePath := path.Join(TmpDir, fileName)
	return compareFileFromGCSBucketAndMntDir(filePathInGcsBucket, mntFilePath, localFilePath, t)
}

func TestMultipleFilesAtSameTime(t *testing.T) {
	concurrentWriteDir := path.Join(setup.MntDir(), DirForConcurrentWrite)
	setup.SetupTestDirectory(DirForConcurrentWrite)

	// Clean up.
	defer operations.RemoveDir(concurrentWriteDir)

	files := []string{FileOne, FileTwo, FileThree}

	var eG errgroup.Group

	// Concurrently write three files.
	for i := range files {
		// Copy the current value of i into a local variable to avoid data races.
		fileIndex := i

		// Thread to write the current file.
		eG.Go(func() error {
			mntFilePath := path.Join(setup.MntDir(), DirForConcurrentWrite, files[fileIndex])
			err := writeFile(mntFilePath, FiveHundredMB)
			if err != nil {
				return fmt.Errorf("WriteError: %v", err)
			}

			return validateFileContents(files[fileIndex], mntFilePath, t)
		})
	}

	// Wait on threads to end.
	err := eG.Wait()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}
