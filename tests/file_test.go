package g_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/enetx/g"
)

func TestFile_Dir_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Get the directory of the file
	result := file.Dir()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestFile_Dir_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the directory of the file is correct
	expectedDir := filepath.Dir(tempFile)
	actualDir := result.Ok().Path().Ok().Std()
	if actualDir != expectedDir {
		t.Errorf("TestFile_Dir_Success: Expected directory %s, got %s", expectedDir, actualDir)
	}
}

func TestFile_Exist_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Check if the file exists
	exists := file.Exist()

	// Check if the existence check is correct
	if !exists {
		t.Errorf("TestFile_Exist_Success: File should exist, but it doesn't.")
	}
}

func TestFile_MimeType_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFileWithData(t, []byte("test content"))
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Get the MIME type of the file
	result := file.MimeType()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestFile_MimeType_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the detected MIME type is correct
	expectedMimeType := "text/plain; charset=utf-8" // MIME type for "test content"
	actualMimeType := result.Ok().Std()
	if actualMimeType != expectedMimeType {
		t.Errorf("TestFile_MimeType_Success: Expected MIME type %s, got %s", expectedMimeType, actualMimeType)
	}
}

func TestFile_Read_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFileWithData(t, []byte("test content"))
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Read the contents of the file
	result := file.Read()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestFile_Read_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the contents of the file are correct
	expectedContent := "test content"
	actualContent := result.Ok().Std()
	if actualContent != expectedContent {
		t.Errorf("TestFile_Read_Success: Expected content %s, got %s", expectedContent, actualContent)
	}
}

func TestFile_IsLink_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a symbolic link to the temporary file
	symlink := tempFile + ".link"
	err := os.Symlink(tempFile, symlink)
	if err != nil {
		t.Fatalf("Failed to create symbolic link: %s", err)
	}
	defer os.Remove(symlink)

	// Create a File instance representing the symbolic link
	file := g.NewFile(g.String(symlink))

	// Check if the file is a symbolic link
	isLink := file.IsLink()

	// Check if the result is correct
	if !isLink {
		t.Errorf("TestFile_IsLink_Success: Expected file to be a symbolic link, but it is not.")
	}
}

func TestFile_IsLink_Failure(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Check if the file is a symbolic link
	isLink := file.IsLink()

	// Check if the result is correct
	if isLink {
		t.Errorf("TestFile_IsLink_Failure: Expected file not to be a symbolic link, but it is.")
	}
}

// createTempFileWithData creates a temporary file with the specified data for testing and returns its path.
func createTempFileWithData(t *testing.T, data []byte) string {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	_, err = tempFile.Write(data)
	if err != nil {
		t.Fatalf("Failed to write data to temporary file: %s", err)
	}

	err = tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temporary file: %s", err)
	}

	return tempFile.Name()
}

func TestFile_Chunks_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write content to the temporary file
	content := "abcdefghijklmnopqrstuvwxyz"
	writeToFile(t, tempFile, content)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Define the chunk size
	chunkSize := g.Int(5)

	// Read the file in chunks
	result := file.Chunks(chunkSize)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_Chunks_Success: Expected Chunks to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Unwrap the Result type to get the underlying iterator
	iterator := result.Ok().Collect()

	// Read chunks from the iterator and verify their content
	expectedChunks := g.Slice[g.String]{"abcde", "fghij", "klmno", "pqrst", "uvwxy", "z"}

	if iterator.Ne(expectedChunks) {
		t.Fatalf(
			"TestFile_Chunks_Success: Expected chunks %v, got %v",
			expectedChunks,
			iterator,
		)
	}
}

func TestFile_Lines_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write content to the temporary file
	content := "line1\nline2\nline3\nline4\nline5\n"
	writeToFile(t, tempFile, content)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Read the file line by line
	result := file.Lines()

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_Lines_Success: Expected Lines to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Unwrap the Result type to get the underlying iterator
	iterator := result.Ok().Collect()

	// Read lines from the iterator and verify their content
	expectedLines := g.Slice[g.String]{"line1", "line2", "line3", "line4", "line5"}

	if iterator.Ne(expectedLines) {
		t.Fatalf(
			"TestFile_Lines_Success: Expected lines %v, got %v",
			expectedLines,
			iterator,
		)
	}
}

func TestFile_Lines_Failure(t *testing.T) {
	// Create a File instance with an invalid path for testing
	file := g.NewFile(g.String("/invalid/path"))

	// Read the file line by line
	result := file.Lines()

	// Check if the result is an error
	if !result.IsErr() {
		t.Fatalf(
			"TestFile_Lines_Failure: Expected Lines to return an error, but got a successful result: %v",
			result.Ok(),
		)
	}
}

func TestFile_Append_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Append content to the file
	content := "appended content"
	result := file.Append(g.String(content))

	defer file.Close()

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_Append_Success: Expected Append to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Read the content of the file to verify the appended content
	fileContent, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("TestFile_Append_Success: Failed to read file content: %v", err)
	}

	// Check if the appended content is present in the file
	if string(fileContent) != content {
		t.Errorf("TestFile_Append_Success: Expected file content to be %s, got %s", content, string(fileContent))
	}
}

func TestFile_Append_Failure(t *testing.T) {
	// Create a File instance with an invalid path for testing
	file := g.NewFile(g.String("/invalid/path"))

	// Append content to the file
	result := file.Append(g.String("appended content"))

	// Check if the result is an error
	if !result.IsErr() {
		t.Fatalf(
			"TestFile_Append_Failure: Expected Append to return an error, but got a successful result: %v",
			result.Ok(),
		)
	}
}

func TestFile_Seek_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write content to the temporary file
	writeToFile(t, tempFile, "test content")

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Seek to the middle of the file
	result := file.Seek(5, io.SeekStart)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_Seek_Success: Expected Seek to return a successful result, but got an error: %v",
			result.Err(),
		)
	}
}

func TestFile_Seek_Failure(t *testing.T) {
	// Create a File instance with an invalid path for testing
	file := g.NewFile(g.String("/invalid/path"))

	// Seek to a position in the file
	result := file.Seek(10, io.SeekStart)

	// Check if the result is an error
	if result.IsOk() {
		t.Fatalf(
			"TestFile_Seek_Failure: Expected Seek to return an error, but got a successful result: %v",
			result.Ok(),
		)
	}

	// Check if the error is of the correct type
	if !os.IsNotExist(result.Err()) {
		t.Fatalf(
			"TestFile_Seek_Failure: Expected error to be of type os.ErrNotExist, but got: %v",
			result.Err(),
		)
	}
}

func TestFile_Rename_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Define the new path for renaming
	newPath := g.NewFile(g.String(tempFile) + "_renamed")
	defer newPath.Remove()

	// Rename the file
	result := file.Rename(newPath.Path().Ok())

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_Rename_Success: Expected Rename to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Verify if the file exists at the new path
	if !newPath.Exist() {
		t.Fatalf("TestFile_Rename_Success: File does not exist at the new path: %s", newPath.Path().Ok().Std())
	}

	// Verify if the original file does not exist
	if file.Exist() {
		t.Fatalf("TestFile_Rename_Success: Original file still exists after renaming: %s", file.Path().Ok().Std())
	}
}

func TestFile_Rename_FileNotExist(t *testing.T) {
	// Create a File instance with an invalid path for testing
	file := g.NewFile(g.String("/invalid/path"))

	// Define the new path for renaming
	newPath := g.NewFile(g.String("/new/path"))

	// Rename the file
	result := file.Rename(newPath.Path().Ok())

	// Check if the result is an error
	if !result.IsErr() {
		t.Fatalf(
			"TestFile_Rename_FileNotExist: Expected Rename to return an error, but got a successful result: %v",
			result.Ok(),
		)
	}

	// Check if the error is of type ErrFileNotExist
	_, ok := result.Err().(*g.ErrFileNotExist)
	if !ok {
		t.Fatalf("TestFile_Rename_FileNotExist: Expected error of type ErrFileNotExist, got: %v", result.Err())
	}
}

func TestFile_OpenFile_ReadLock(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with read-lock
	result := file.OpenFile(os.O_RDONLY, 0644)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_OpenFile_ReadLock: Expected OpenFile to return a successful result, but got an error: %v",
			result.Err(),
		)
	}
}

func TestFile_OpenFile_WriteLock(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with write-lock
	result := file.OpenFile(os.O_WRONLY, 0644)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_OpenFile_WriteLock: Expected OpenFile to return a successful result, but got an error: %v",
			result.Err(),
		)
	}
}

func TestFile_OpenFile_Truncate(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with truncation flag
	result := file.OpenFile(os.O_WRONLY|os.O_TRUNC, 0644)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_OpenFile_Truncate: Expected OpenFile to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Verify if the file is truncated
	fileStat, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("TestFile_OpenFile_Truncate: Failed to get file stat: %v", err)
	}
	if fileStat.Size() != 0 {
		t.Fatalf("TestFile_OpenFile_Truncate: File is not truncated")
	}
}

func TestFile_Chmod(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Change the mode of the file
	result := file.Chmod(0644)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf("TestFile_Chmod: Expected Chmod to return a successful result, but got an error: %v", result.Err())
	}

	// Verify if the mode of the file is changed
	fileStat, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("TestFile_Chmod: Failed to get file stat: %v", err)
	}
	expectedMode := os.FileMode(0644)
	if fileStat.Mode() != expectedMode {
		t.Fatalf("TestFile_Chmod: Expected file mode to be %v, but got %v", expectedMode, fileStat.Mode())
	}
}

func TestFile_Chown(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Change the owner of the file
	result := file.Chown(os.Getuid(), os.Getgid())

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf("TestFile_Chown: Expected Chown to return a successful result, but got an error: %v", result.Err())
	}

	// Verify if the owner of the file is changed
	fileStat, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("TestFile_Chown: Failed to get file stat: %v", err)
	}
	expectedUID := os.Getuid()
	expectedGID := os.Getgid()
	statT, ok := fileStat.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatalf("TestFile_Chown: Failed to get file Sys info")
	}
	if int(statT.Uid) != expectedUID || int(statT.Gid) != expectedGID {
		t.Fatalf(
			"TestFile_Chown: Expected file owner to be UID: %d, GID: %d, but got UID: %d, GID: %d",
			expectedUID,
			expectedGID,
			statT.Uid,
			statT.Gid,
		)
	}
}

func TestFile_WriteFromReader(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Prepare data to write
	testData := "Hello, World!"
	reader := bytes.NewBufferString(testData)

	// Write data from the reader into the file
	result := file.WriteFromReader(reader)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf(
			"TestFile_WriteFromReader: Expected WriteFromReader to return a successful result, but got an error: %v",
			result.Err(),
		)
	}

	// Read the content of the file to verify if the data is written correctly
	contentResult := file.Read()
	if contentResult.IsErr() {
		t.Fatalf("TestFile_WriteFromReader: Failed to read file content: %v", contentResult.Err())
	}

	// Verify if the content of the file matches the test data
	if contentResult.Ok().Std() != testData {
		t.Fatalf(
			"TestFile_WriteFromReader: Expected file content to be '%s', but got '%s'",
			testData,
			contentResult.Ok().Std(),
		)
	}
}

func TestFile_Write(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Prepare data to write
	testData := "Hello, World!"

	// Write data into the file
	result := file.Write(g.String(testData))

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf("TestFile_Write: Expected Write to return a successful result, but got an error: %v", result.Err())
	}

	// Read the content of the file to verify if the data is written correctly
	contentResult := file.Read()
	if contentResult.IsErr() {
		t.Fatalf("TestFile_Write: Failed to read file content: %v", contentResult.Err())
	}

	// Verify if the content of the file matches the test data
	if contentResult.Ok().Std() != testData {
		t.Fatalf("TestFile_Write: Expected file content to be '%s', but got '%s'", testData, contentResult.Ok().Std())
	}
}

func TestFile_Ext(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Extract the file extension
	extension := file.Ext().Std()

	// Expected file extension (assuming the temporary file has an extension)
	expectedExtension := ".txt"

	// Check if the extracted extension matches the expected extension
	if extension != expectedExtension {
		t.Fatalf("TestFile_Ext: Expected extension to be '%s', but got '%s'", expectedExtension, extension)
	}
}

func TestFile_Copy(t *testing.T) {
	// Create a temporary source file for testing
	srcFile := createTempFile(t)
	defer os.Remove(srcFile)

	// Create a temporary destination file for testing
	destFile := createTempFile(t)
	defer os.Remove(destFile)

	// Create a File instance representing the source file
	src := g.NewFile(g.String(srcFile))

	// Copy the source file to the destination file
	dest := g.NewFile(g.String(destFile))
	result := src.Copy(dest.Path().Ok())

	// Check if the copy operation was successful
	if result.IsErr() {
		t.Fatalf("TestFile_Copy: Failed to copy file: %s", result.Err())
	}

	// Verify that the destination file exists
	if !dest.Exist() {
		t.Fatalf("TestFile_Copy: Destination file does not exist after copy")
	}
}

func TestFile_Split(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Split the file path into its directory and file components
	dir, fileName := file.Split()

	// Check if the directory and file components are correct
	if dir == nil || dir.Path().Ok().Std() != filepath.Dir(tempFile) {
		t.Errorf("TestFile_Split: Incorrect directory component")
	}

	if fileName == nil || fileName.Name().Std() != filepath.Base(tempFile) {
		t.Errorf("TestFile_Split: Incorrect file name component")
	}
}

func TestFile_Move(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := g.NewFile(g.String(tempFile))

	// Rename the file
	newpath := tempFile + ".new"
	renamedFile := file.Move(g.String(newpath))

	// Check if the file has been successfully renamed
	if renamedFile.IsErr() {
		t.Errorf("TestFile_Rename: Failed to rename file: %v", renamedFile.Err())
	}

	defer renamedFile.Ok().Remove()

	// Verify that the old file does not exist
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Errorf("TestFile_Rename: Old file still exists after renaming")
	}

	// Verify that the new file exists
	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		t.Errorf("TestFile_Rename: New file does not exist after renaming")
	}
}

// writeToFile writes content to the specified file.
func writeToFile(t *testing.T, filename, content string) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open file %s for writing: %s", filename, err)
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(content))
	if err != nil {
		t.Fatalf("Failed to write content to file %s: %s", filename, err)
	}
}

// createTempFile creates a temporary file for testing and returns its path.
func createTempFile(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	defer tempFile.Close()

	return tempFile.Name()
}
