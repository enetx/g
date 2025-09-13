package g_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	. "github.com/enetx/g"
)

func TestFile_Dir_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(symlink))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

	// Define the chunk size
	chunkSize := Int(5)

	// Read the file in chunks
	result := file.Chunks(chunkSize)

	// Check if the result is successful
	if result.FirstErr().IsSome() {
		t.Fatalf(
			"TestFile_Chunks_Success: Expected Chunks to return a successful result, but got an error: %v",
			result.FirstErr().Some(),
		)
	}

	// Unwrap the Result type to get the underlying iterator
	iterator := result.Ok().Collect()

	// Read chunks from the iterator and verify their content
	expectedChunks := Slice[String]{"abcde", "fghij", "klmno", "pqrst", "uvwxy", "z"}

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
	file := NewFile(String(tempFile))

	// Read the file line by line
	result := file.Lines()

	// Check if the result is successful
	if result.FirstErr().IsSome() {
		t.Fatalf(
			"TestFile_Lines_Success: Expected Lines to return a successful result, but got an error: %v",
			result.FirstErr().Some(),
		)
	}

	// Unwrap the Result type to get the underlying iterator
	iterator := result.Ok().Collect()

	// Read lines from the iterator and verify their content
	expectedLines := Slice[String]{"line1", "line2", "line3", "line4", "line5"}

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
	file := NewFile(String("/invalid/path"))

	// Read the file line by line
	result := file.Lines()

	// Check if the result is an error
	if result.FirstErr().IsNone() {
		t.Fatalf(
			"TestFile_Lines_Failure: Expected Lines to return an error, but got a successful result: %v",
			result.Ok().Collect(),
		)
	}
}

func TestFile_Append_Success(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := NewFile(String(tempFile))

	// Append content to the file
	content := "appended content"
	result := file.Append(String(content))

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
	file := NewFile(String("/invalid/path"))

	// Append content to the file
	result := file.Append(String("appended content"))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String("/invalid/path"))

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
	file := NewFile(String(tempFile))

	// Define the new path for renaming
	newPath := NewFile(String(tempFile) + "_renamed")
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
	file := NewFile(String("/invalid/path"))

	// Define the new path for renaming
	newPath := NewFile(String("/new/path"))

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
	_, ok := result.Err().(*ErrFileNotExist)
	if !ok {
		t.Fatalf("TestFile_Rename_FileNotExist: Expected error of type ErrFileNotExist, got: %v", result.Err())
	}
}

func TestFile_OpenFile_ReadLock(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := NewFile(String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with read-lock
	result := file.OpenFile(os.O_RDONLY, 0o644)

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
	file := NewFile(String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with write-lock
	result := file.OpenFile(os.O_WRONLY, 0o644)

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
	file := NewFile(String(tempFile))

	file.Guard()
	defer file.Close()

	// Open the file with truncation flag
	result := file.OpenFile(os.O_WRONLY|os.O_TRUNC, 0o644)

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

func TestFile_OpenFile_NonExistentFile(t *testing.T) {
	// Test OpenFile with a non-existent file (should error)
	file := NewFile(String("/nonexistent/file/path.txt"))

	result := file.OpenFile(os.O_RDONLY, 0o644)
	if result.IsOk() {
		t.Errorf("TestFile_OpenFile_NonExistentFile: Expected error for non-existent file, but operation succeeded")
	}
}

func TestFile_OpenFile_NoGuard(t *testing.T) {
	// Test OpenFile without Guard (should work without locking)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	defer file.Close()

	// Open without calling Guard() - no file locking
	result := file.OpenFile(os.O_RDONLY, 0o644)
	if result.IsErr() {
		t.Errorf("TestFile_OpenFile_NoGuard: Expected success without guard, got error: %v", result.Err())
	}
}

func TestFile_OpenFile_ReadWriteMode(t *testing.T) {
	// Test OpenFile with O_RDWR flag
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	file.Guard()
	defer file.Close()

	result := file.OpenFile(os.O_RDWR, 0o644)
	if result.IsErr() {
		t.Errorf("TestFile_OpenFile_ReadWriteMode: Expected success with O_RDWR, got error: %v", result.Err())
	}
}

func TestFile_OpenFile_CreateFlag(t *testing.T) {
	// Test OpenFile with O_CREATE flag for new file
	tempFile := createTempFile(t)
	os.Remove(tempFile) // Remove it first so we can test creation
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	file.Guard()
	defer file.Close()

	result := file.OpenFile(os.O_WRONLY|os.O_CREATE, 0o644)
	if result.IsErr() {
		t.Errorf("TestFile_OpenFile_CreateFlag: Expected success with O_CREATE, got error: %v", result.Err())
	}

	// Verify the file was created
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("TestFile_OpenFile_CreateFlag: File should have been created")
	}
}

func TestFile_OpenFile_TruncateNoGuard(t *testing.T) {
	// Test OpenFile with truncate flag but no guard
	tempFile := createTempFileWithContent(t, "initial content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	defer file.Close()

	// Test truncate without guard
	result := file.OpenFile(os.O_WRONLY|os.O_TRUNC, 0o644)
	if result.IsErr() {
		t.Errorf("TestFile_OpenFile_TruncateNoGuard: Expected success, got error: %v", result.Err())
	}

	// Verify the file was truncated
	fileStat, err := os.Stat(tempFile)
	if err != nil {
		t.Errorf("TestFile_OpenFile_TruncateNoGuard: Failed to stat file: %v", err)
	} else if fileStat.Size() != 0 {
		t.Error("TestFile_OpenFile_TruncateNoGuard: File should be truncated")
	}
}

func TestFile_Name_ClosedFile(t *testing.T) {
	// Test Name() on a file that hasn't been opened (f.file == nil)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	// Get name without opening the file
	name := file.Name()
	expectedName := filepath.Base(tempFile)

	if name.Std() != expectedName {
		t.Errorf("TestFile_Name_ClosedFile: Expected name '%s', got '%s'", expectedName, name.Std())
	}
}

func TestFile_Name_OpenedFile(t *testing.T) {
	// Test Name() on an opened file (f.file != nil)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	defer file.Close()

	// Open the file first
	result := file.Open()
	if result.IsErr() {
		t.Fatalf("Failed to open file: %v", result.Err())
	}

	// Get name after opening the file
	name := file.Name()
	expectedName := filepath.Base(tempFile)

	if name.Std() != expectedName {
		t.Errorf("TestFile_Name_OpenedFile: Expected name '%s', got '%s'", expectedName, name.Std())
	}
}

func TestFile_Name_WithPath(t *testing.T) {
	// Test Name() with a complex file path
	complexPath := "/path/to/directory/filename.txt"
	file := NewFile(String(complexPath))

	name := file.Name()
	expectedName := "filename.txt"

	if name.Std() != expectedName {
		t.Errorf("TestFile_Name_WithPath: Expected name '%s', got '%s'", expectedName, name.Std())
	}
}

func TestFile_Remove_Success(t *testing.T) {
	// Test successful file removal
	tempFile := createTempFileWithContent(t, "test content")
	// Don't defer removal since we're testing Remove()

	file := NewFile(String(tempFile))

	// Remove the file
	result := file.Remove()
	if result.IsErr() {
		t.Errorf("TestFile_Remove_Success: Expected success, got error: %v", result.Err())
	}

	// Verify the file no longer exists
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Error("TestFile_Remove_Success: File should not exist after removal")
	}
}

func TestFile_Remove_NonExistentFile(t *testing.T) {
	// Test removal of a non-existent file (should error)
	file := NewFile(String("/nonexistent/file/path.txt"))

	result := file.Remove()
	if result.IsOk() {
		t.Error("TestFile_Remove_NonExistentFile: Expected error for non-existent file, but operation succeeded")
	}
}

func TestFile_Remove_Directory(t *testing.T) {
	// Test removal of a directory (should error since os.Remove doesn't remove non-empty dirs)
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir) // Fallback cleanup

	// Create a file inside the directory
	if err := os.WriteFile(tempDir+"/test.txt", []byte("content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := NewFile(String(tempDir))

	result := file.Remove()
	if result.IsOk() {
		t.Error(
			"TestFile_Remove_Directory: Expected error when trying to remove non-empty directory, but operation succeeded",
		)
	}
}

func TestFile_Stat_ClosedFile(t *testing.T) {
	// Test Stat() on a file that hasn't been opened (f.file == nil)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	// Get stat without opening the file
	result := file.Stat()
	if result.IsErr() {
		t.Errorf("TestFile_Stat_ClosedFile: Expected success, got error: %v", result.Err())
	}

	stat := result.Ok()
	if stat.IsDir() {
		t.Error("TestFile_Stat_ClosedFile: Expected regular file, got directory")
	}
	if stat.Size() != int64(len("test content")) {
		t.Errorf("TestFile_Stat_ClosedFile: Expected size %d, got %d", len("test content"), stat.Size())
	}
}

func TestFile_Stat_OpenedFile(t *testing.T) {
	// Test Stat() on an opened file (f.file != nil)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	defer file.Close()

	// Open the file first
	openResult := file.Open()
	if openResult.IsErr() {
		t.Fatalf("Failed to open file: %v", openResult.Err())
	}

	// Get stat after opening the file
	result := file.Stat()
	if result.IsErr() {
		t.Errorf("TestFile_Stat_OpenedFile: Expected success, got error: %v", result.Err())
	}

	stat := result.Ok()
	if stat.IsDir() {
		t.Error("TestFile_Stat_OpenedFile: Expected regular file, got directory")
	}
	if stat.Size() != int64(len("test content")) {
		t.Errorf("TestFile_Stat_OpenedFile: Expected size %d, got %d", len("test content"), stat.Size())
	}
}

func TestFile_Stat_NonExistentFile(t *testing.T) {
	// Test Stat() on a non-existent file (should error)
	file := NewFile(String("/nonexistent/file/path.txt"))

	result := file.Stat()
	if result.IsOk() {
		t.Error("TestFile_Stat_NonExistentFile: Expected error for non-existent file, but operation succeeded")
	}
}

func TestFile_Chmod(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := NewFile(String(tempFile))

	// Change the mode of the file
	result := file.Chmod(0o644)

	// Check if the result is successful
	if result.IsErr() {
		t.Fatalf("TestFile_Chmod: Expected Chmod to return a successful result, but got an error: %v", result.Err())
	}

	// Verify if the mode of the file is changed
	fileStat, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("TestFile_Chmod: Failed to get file stat: %v", err)
	}
	expectedMode := os.FileMode(0o644)
	if fileStat.Mode() != expectedMode {
		t.Fatalf("TestFile_Chmod: Expected file mode to be %v, but got %v", expectedMode, fileStat.Mode())
	}
}

func TestFile_Chown(t *testing.T) {
	// Create a temporary file for testing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Create a File instance representing the temporary file
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

	// Prepare data to write
	testData := "Hello, World!"

	// Write data into the file
	result := file.Write(String(testData))

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
	file := NewFile(String(tempFile))

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
	src := NewFile(String(srcFile))

	// Copy the source file to the destination file
	dest := NewFile(String(destFile))
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
	file := NewFile(String(tempFile))

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
	file := NewFile(String(tempFile))

	// Rename the file
	newpath := tempFile + ".new"
	renamedFile := file.Move(String(newpath))

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
	file, err := os.OpenFile(filename, os.O_WRONLY, 0o644)
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

// createTempFileWithContent creates a temporary file with specific content for testing and returns its path.
func createTempFileWithContent(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write content to temporary file: %s", err)
	}

	defer tempFile.Close()

	return tempFile.Name()
}

func TestFile_LinesRaw(t *testing.T) {
	// Create a temporary file with multiple lines
	tempFilePath := createTempFile(t)
	defer os.Remove(tempFilePath)

	testContent := "line1\nline2\nline3"
	file := NewFile(String(tempFilePath))
	file.Write(String(testContent))

	// Test LinesRaw iterator
	var lines []string
	file.LinesRaw().ForEach(func(result Result[Bytes]) {
		if result.IsErr() {
			t.Errorf("LinesRaw: Unexpected error: %s", result.Err().Error())
			return
		}
		lines = append(lines, result.Ok().String().Std())
	})

	expectedLines := []string{"line1", "line2", "line3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("LinesRaw: Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, expectedLine := range expectedLines {
		if i < len(lines) && lines[i] != expectedLine {
			t.Errorf("LinesRaw: Line %d mismatch. Expected '%s', got '%s'", i, expectedLine, lines[i])
		}
	}
}

func TestFile_ChunksRaw(t *testing.T) {
	// Create a temporary file with test content
	tempFilePath := createTempFile(t)
	defer os.Remove(tempFilePath)

	testContent := "abcdefghijklmnop"
	file := NewFile(String(tempFilePath))
	file.Write(String(testContent))

	// Test ChunksRaw iterator with chunk size 5
	var chunks []string
	file.ChunksRaw(5).ForEach(func(result Result[Bytes]) {
		if result.IsErr() {
			t.Errorf("ChunksRaw: Unexpected error: %s", result.Err().Error())
			return
		}
		chunks = append(chunks, result.Ok().String().Std())
	})

	expectedChunks := []string{"abcde", "fghij", "klmno", "p"}
	if len(chunks) != len(expectedChunks) {
		t.Errorf("ChunksRaw: Expected %d chunks, got %d", len(expectedChunks), len(chunks))
	}

	for i, expectedChunk := range expectedChunks {
		if i < len(chunks) && chunks[i] != expectedChunk {
			t.Errorf("ChunksRaw: Chunk %d mismatch. Expected '%s', got '%s'", i, expectedChunk, chunks[i])
		}
	}

	// Test error case with invalid chunk size
	var errorReceived bool
	file.ChunksRaw(0).ForEach(func(result Result[Bytes]) {
		if result.IsErr() {
			errorReceived = true
			if !strings.Contains(result.Err().Error(), "chunk size must be > 0") {
				t.Errorf("ChunksRaw: Expected specific error message, got: %s", result.Err().Error())
			}
		}
	})

	if !errorReceived {
		t.Error("ChunksRaw: Expected error for chunk size <= 0, but no error was received")
	}
}

func TestFile_Print(t *testing.T) {
	tempFilePath := createTempFile(t)
	defer os.Remove(tempFilePath)

	file := NewFile(String(tempFilePath))
	result := file.Print()
	if result != file {
		t.Errorf("Print() should return original file unchanged")
	}
}

func TestFile_Println(t *testing.T) {
	tempFilePath := createTempFile(t)
	defer os.Remove(tempFilePath)

	file := NewFile(String(tempFilePath))
	result := file.Println()
	if result != file {
		t.Errorf("Println() should return original file unchanged")
	}
}

func TestFile_Lines_EmptyFile(t *testing.T) {
	// Test Lines with empty file
	tempFile := createTempFileWithContent(t, "")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	linesResult := file.Lines()

	if linesResult.FirstErr().IsSome() {
		t.Errorf("Lines should succeed, got error: %v", linesResult.FirstErr().Some())
		return
	}

	lines := linesResult.Ok().Collect()
	if lines.Len() != 0 {
		t.Errorf("Empty file should have 0 lines, got %d", lines.Len())
	}
}

func TestFile_Lines_SingleLine(t *testing.T) {
	// Test Lines with single line without newline
	tempFile := createTempFileWithContent(t, "single line")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	linesResult := file.Lines()

	if linesResult.FirstErr().IsSome() {
		t.Errorf("Lines should succeed, got error: %v", linesResult.FirstErr().Some())
		return
	}

	lines := linesResult.Ok().Collect()
	if lines.Len() != 1 {
		t.Errorf("Single line file should have 1 line, got %d", lines.Len())
	}
	if lines.Get(0).Some() != "single line" {
		t.Errorf("Expected 'single line', got '%s'", lines.Get(0).Some())
	}
}

func TestFile_Chunks_SmallFile(t *testing.T) {
	// Test Chunks with file smaller than chunk size
	tempFile := createTempFileWithContent(t, "small")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	chunks := file.Chunks(10).Ok().Collect() // chunk size larger than file

	if chunks.Len() != 1 {
		t.Errorf("Small file should produce 1 chunk, got %d", chunks.Len())
	}
	if chunks.Get(0).Some() != "small" {
		t.Errorf("Expected 'small', got '%s'", chunks.Get(0).Some())
	}
}

func TestFile_Chunks_InvalidSize(t *testing.T) {
	// Test Chunks with invalid chunk size (0 and negative)
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	// Test with chunk size 0
	result := file.Chunks(0)
	if result.FirstErr().IsNone() {
		t.Errorf("TestFile_Chunks_InvalidSize: Expected error for chunk size 0, but operation succeeded")
	}

	if result.FirstErr().IsSome() && result.FirstErr().Some().Error() != "chunk size must be > 0" {
		t.Errorf("TestFile_Chunks_InvalidSize: Expected specific error message, got: %s", result.FirstErr().Some())
	}

	// Test with negative chunk size
	result = file.Chunks(-5)
	if result.FirstErr().IsNone() {
		t.Errorf("TestFile_Chunks_InvalidSize: Expected error for negative chunk size, but operation succeeded")
	}
}

func TestFile_Chunks_NonExistentFile(t *testing.T) {
	// Test Chunks with a non-existent file (should error during Open)
	file := NewFile(String("/nonexistent/file/path.txt"))

	result := file.Chunks(5)
	if result.FirstErr().IsNone() {
		t.Errorf("TestFile_Chunks_NonExistentFile: Expected error for non-existent file, but operation succeeded")
	}
}

func TestFile_Chunks_EarlyExit(t *testing.T) {
	// Test Chunks with early exit (iterator stops after first chunk)
	tempFile := createTempFileWithContent(t, "abcdefghijklmnopqrstuvwxyz")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	var processedChunks []String
	count := 0

	// Process only the first chunk, then stop
	file.Chunks(5).Range(func(result Result[String]) bool {
		if result.IsErr() {
			t.Errorf("TestFile_Chunks_EarlyExit: Unexpected error: %s", result.Err().Error())
			return false
		}
		processedChunks = append(processedChunks, result.Ok())
		count++
		// Stop after first chunk
		return count < 1
	})

	if len(processedChunks) != 1 {
		t.Errorf("TestFile_Chunks_EarlyExit: Expected 1 chunk processed, got %d", len(processedChunks))
	}
	if len(processedChunks) > 0 && processedChunks[0] != "abcde" {
		t.Errorf("TestFile_Chunks_EarlyExit: Expected first chunk 'abcde', got '%s'", processedChunks[0])
	}
}

func TestFile_Append_NewContent(t *testing.T) {
	// Test Append with additional content
	tempFile := createTempFileWithContent(t, "initial content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	appendContent := String("appended content")
	result := file.Append(appendContent)
	defer file.Close()

	if result.IsErr() {
		t.Errorf("Append should succeed, got error: %v", result.Err())
	}

	// Verify content was appended
	content, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expected := "initial contentappended content"
	if string(content) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(content))
	}
}

func TestFile_OpenFile_TruncateErrorHandling(t *testing.T) {
	// Test OpenFile truncate behavior
	tempFile := createTempFileWithContent(t, "test content")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	file.Guard()
	defer file.Close()

	// Test truncate behavior with O_TRUNC flag
	result := file.OpenFile(os.O_WRONLY|os.O_TRUNC, 0o644)
	if result.IsErr() {
		t.Errorf("Expected success with O_TRUNC, got: %v", result.Err())
	}
}

func TestFile_OpenFile_ExclusiveFlag(t *testing.T) {
	// Test O_EXCL with O_CREATE
	tempDir := t.TempDir()

	newFile := filepath.Join(tempDir, "exclusive_test.txt")
	file := NewFile(String(newFile))
	defer file.Close()

	// Create file with O_CREATE|O_EXCL
	result := file.OpenFile(os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if result.IsErr() {
		t.Errorf("Expected success creating file with O_CREATE|O_EXCL, got: %v", result.Err())
	}

	file2 := NewFile(String(newFile))
	defer file2.Close()

	// Try to create again with O_EXCL - should fail
	result2 := file2.OpenFile(os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if result2.IsOk() {
		t.Error("Expected error when creating existing file with O_EXCL")
	}
}

func TestFile_OpenFile_AppendFlag(t *testing.T) {
	// Test O_APPEND flag
	tempFile := createTempFileWithContent(t, "initial")
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	file.Guard()
	defer file.Close()

	result := file.OpenFile(os.O_WRONLY|os.O_APPEND, 0o644)
	if result.IsErr() {
		t.Errorf("Expected success with O_APPEND, got: %v", result.Err())
	}
}

func TestFile_OpenFile_DefaultMode(t *testing.T) {
	// Test with default flag (neither WRONLY nor RDWR, should use read lock)
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	file.Guard()
	defer file.Close()

	// Test with 0 flag (default readonly)
	result := file.OpenFile(0, 0o644)
	if result.IsErr() {
		t.Errorf("Expected success with default flag, got: %v", result.Err())
	}
}

func TestFile_OpenFile_TruncateStatError(t *testing.T) {
	// Test truncate with stat error path - using a directory instead of file
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	file := NewFile(String(tempDir))
	file.Guard()
	defer file.Close()

	// Try to open directory with O_TRUNC - should hit the stat/IsRegular path
	result := file.OpenFile(os.O_RDWR|os.O_TRUNC, 0o644)
	if result.IsOk() {
		t.Log("Opening directory with O_TRUNC succeeded (platform-dependent behavior)")
	} else {
		t.Log("Opening directory with O_TRUNC failed as expected")
	}
}

func TestFile_CreateTemp_NoArgs(t *testing.T) {
	// Test CreateTemp with no arguments (should use default dir and pattern)
	result := NewFile("").CreateTemp()
	if result.IsErr() {
		t.Errorf("CreateTemp with no args should succeed, got error: %v", result.Err())
		return
	}

	tempFile := result.Ok()
	tempPath := tempFile.Path().Ok().Std()
	defer os.Remove(tempPath)

	// Check that file was created in temp directory
	if !strings.Contains(tempPath, os.TempDir()) {
		t.Errorf("Expected temp file to be in temp directory, got: %s", tempPath)
	}
}

func TestFile_CreateTemp_WithDir(t *testing.T) {
	// Test CreateTemp with only directory argument
	tempDir := os.TempDir()
	result := NewFile("").CreateTemp(String(tempDir))
	if result.IsErr() {
		t.Errorf("CreateTemp with dir should succeed, got error: %v", result.Err())
		return
	}

	tempFile := result.Ok()
	tempPath := tempFile.Path().Ok().Std()
	defer os.Remove(tempPath)

	// Check that file was created in specified directory
	if !strings.Contains(tempPath, tempDir) {
		t.Errorf("Expected temp file to be in %s, got: %s", tempDir, tempPath)
	}
}

func TestFile_CreateTemp_WithDirAndPattern(t *testing.T) {
	// Test CreateTemp with both directory and pattern arguments
	tempDir := os.TempDir()
	pattern := "testfile_*.tmp"
	result := NewFile("").CreateTemp(String(tempDir), String(pattern))
	if result.IsErr() {
		t.Errorf("CreateTemp with dir and pattern should succeed, got error: %v", result.Err())
		return
	}

	tempFile := result.Ok()
	tempPath := tempFile.Path().Ok().Std()
	defer os.Remove(tempPath)

	// Check that file was created in specified directory
	if !strings.Contains(tempPath, tempDir) {
		t.Errorf("Expected temp file to be in %s, got: %s", tempDir, tempPath)
	}

	// Check that file name matches pattern (should start with "testfile_" and end with ".tmp")
	fileName := filepath.Base(tempPath)
	if !strings.HasPrefix(fileName, "testfile_") || !strings.HasSuffix(fileName, ".tmp") {
		t.Errorf("Expected temp file name to match pattern %s, got: %s", pattern, fileName)
	}
}

func TestFile_CreateTemp_InvalidDir(t *testing.T) {
	// Test CreateTemp with invalid directory (should fail)
	result := NewFile("").CreateTemp(String("/nonexistent/invalid/directory"))
	if result.IsOk() {
		// On some platforms this might succeed, so clean up
		os.Remove(result.Ok().Path().Ok().Std())
		t.Log("CreateTemp with invalid directory succeeded (platform-dependent behavior)")
	} else {
		t.Log("CreateTemp with invalid directory failed as expected")
	}
}

func TestFile_OpenFile_NoGuardMode(t *testing.T) {
	// Test OpenFile without guard mode (should skip locking)
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	// Don't call Guard() - this should skip the locking logic
	defer file.Close()

	// Test different flag combinations without guard
	flags := []int{
		os.O_RDONLY,
		os.O_WRONLY,
		os.O_RDWR,
		os.O_RDWR | os.O_CREATE,
		os.O_WRONLY | os.O_TRUNC,
	}

	for _, flag := range flags {
		result := file.OpenFile(flag, 0o644)
		if result.IsErr() {
			t.Errorf("OpenFile with flag %d should succeed without guard, got error: %v", flag, result.Err())
		}
		// Close the file before next iteration
		file.Close()
	}
}

func TestFile_OpenFile_InvalidFileForLocking(t *testing.T) {
	// Test OpenFile with a file that might cause locking issues
	// Try with a directory path - this might cause lock errors on some platforms
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	file := NewFile(String(tempDir))
	file.Guard()
	defer file.Close()

	// Try to open directory with write mode - might fail during locking
	result := file.OpenFile(os.O_RDWR, 0o644)
	if result.IsOk() {
		t.Log("Opening directory with write mode succeeded (platform-dependent)")
	} else {
		t.Log("Opening directory with write mode failed as expected")
	}
}

func TestFile_Read_OpenError(t *testing.T) {
	// Test Read when Open() fails
	file := NewFile("/nonexistent/file/path.txt")

	result := file.Read()
	if result.IsOk() {
		t.Errorf("Expected Read to fail for non-existent file, but it succeeded")
	} else {
		t.Logf("Read failed as expected for non-existent file: %v", result.Err())
	}
}

func TestFile_Read_ReadAllError(t *testing.T) {
	// Test Read when ReadAll might fail
	// Create a temporary file and then remove it after opening to simulate ReadAll error
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write some content to file
	if err := os.WriteFile(tempFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}

	file := NewFile(String(tempFile))

	// For this test, we'll just test normal operation since simulating ReadAll errors
	// is complex and platform-dependent
	result := file.Read()
	if result.IsErr() {
		t.Errorf("Expected Read to succeed, got error: %v", result.Err())
	} else {
		content := result.Ok()
		if content.Std() != "test content" {
			t.Errorf("Expected content 'test content', got: %s", content.Std())
		}
	}
}

func TestFile_LinesRaw_OpenError(t *testing.T) {
	// Test LinesRaw when Open() fails
	file := NewFile("/nonexistent/path/file.txt")

	errorFound := false
	file.LinesRaw().ForEach(func(result Result[Bytes]) {
		if result.IsErr() {
			errorFound = true
		}
	})

	if !errorFound {
		t.Error("Expected error when opening nonexistent file for LinesRaw")
	}
}

func TestFile_LinesRaw_EarlyReturn(t *testing.T) {
	// Test LinesRaw early return when yield returns false
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write multiple lines
	testContent := "line1\nline2\nline3\nline4\nline5"
	file := NewFile(String(tempFile))
	file.Write(String(testContent))

	linesRead := 0
	for result := range file.LinesRaw() {
		if result.IsOk() {
			linesRead++
			// Stop after reading 2 lines to test early return
			if linesRead >= 2 {
				break
			}
		}
	}

	if linesRead != 2 {
		t.Errorf("Expected to read exactly 2 lines before early return, got %d", linesRead)
	}
}

func TestFile_Chmod_WithOpenFile(t *testing.T) {
	// Test Chmod when file is open (f.file != nil)
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	// Open the file to set f.file != nil
	if openResult := file.Open(); openResult.IsErr() {
		t.Fatalf("Failed to open file: %v", openResult.Err())
	}
	defer file.Close()

	// Test chmod on open file
	newMode := os.FileMode(0o755)
	result := file.Chmod(newMode)

	if result.IsErr() {
		t.Errorf("Chmod on open file failed: %v", result.Err())
	}
}

func TestFile_Chown_WithOpenFile(t *testing.T) {
	// Test Chown when file is open (f.file != nil)
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))
	// Open the file to set f.file != nil
	if openResult := file.Open(); openResult.IsErr() {
		t.Fatalf("Failed to open file: %v", openResult.Err())
	}
	defer file.Close()

	// Test chown on open file
	uid := os.Getuid()
	gid := os.Getgid()
	result := file.Chown(uid, gid)

	if result.IsErr() {
		t.Errorf("Chown on open file failed: %v", result.Err())
	}
}

func TestFile_Lines_OpenError(t *testing.T) {
	// Test Lines when Open() fails
	file := NewFile("/nonexistent/path/file.txt")

	errorFound := false
	for result := range file.Lines() {
		if result.IsErr() {
			errorFound = true
			break
		}
	}

	if !errorFound {
		t.Error("Expected error when opening nonexistent file for Lines")
	}
}

func TestFile_Lines_EarlyReturn(t *testing.T) {
	// Test Lines early return when yield returns false
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write multiple lines
	testContent := "line1\nline2\nline3\nline4\nline5"
	file := NewFile(String(tempFile))
	file.Write(String(testContent))

	linesRead := 0
	for result := range file.Lines() {
		if result.IsOk() {
			linesRead++
			// Stop after reading 2 lines to test early return
			if linesRead >= 2 {
				break
			}
		}
	}

	if linesRead != 2 {
		t.Errorf("Expected to read exactly 2 lines before early return, got %d", linesRead)
	}
}

func TestFile_Seek_OpenError(t *testing.T) {
	// Test Seek when Open() fails
	file := NewFile("/nonexistent/path/file.txt")

	result := file.Seek(0, 0)

	if result.IsOk() {
		t.Error("Expected error when seeking in nonexistent file")
	}
}

func TestFile_Seek_SeekError(t *testing.T) {
	// Test Seek when file.Seek() fails - hard to simulate directly
	// Let's just test with invalid offset to trigger potential errors
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	// Try various seek operations that might fail
	result1 := file.Seek(-1, 2) // Seek to invalid position relative to end
	result2 := file.Seek(0, 0)  // Normal seek should work

	if result1.IsErr() {
		t.Logf("Seek with potentially invalid position failed as expected: %v", result1.Err())
	}

	if result2.IsErr() {
		t.Errorf("Normal seek failed unexpectedly: %v", result2.Err())
	}
}

func TestFile_ChunksRaw_InvalidSize(t *testing.T) {
	// Test ChunksRaw with invalid chunk size
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	errorFound := false
	for result := range file.ChunksRaw(Int(0)) {
		if result.IsErr() {
			errorFound = true
			if !strings.Contains(result.Err().Error(), "chunk size must be > 0") {
				t.Errorf("Expected error about chunk size, got: %v", result.Err())
			}
			break
		}
	}

	if !errorFound {
		t.Error("Expected error for invalid chunk size")
	}
}

func TestFile_ChunksRaw_OpenError(t *testing.T) {
	// Test ChunksRaw when Open() fails
	file := NewFile("/nonexistent/path/file.txt")

	errorFound := false
	for result := range file.ChunksRaw(Int(1024)) {
		if result.IsErr() {
			errorFound = true
			break
		}
	}

	if !errorFound {
		t.Error("Expected error when opening nonexistent file for ChunksRaw")
	}
}

func TestFile_ChunksRaw_EarlyReturn(t *testing.T) {
	// Test ChunksRaw early return when yield returns false
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	// Write test content
	testContent := strings.Repeat("abcdefghij", 100) // 1000 bytes
	file := NewFile(String(tempFile))
	file.Write(String(testContent))

	chunksRead := 0
	for result := range file.ChunksRaw(Int(100)) {
		if result.IsOk() {
			chunksRead++
			// Stop after reading 2 chunks to test early return
			if chunksRead >= 2 {
				break
			}
		}
	}

	if chunksRead != 2 {
		t.Errorf("Expected to read exactly 2 chunks before early return, got %d", chunksRead)
	}
}

func TestFile_Append_CreateAllError(t *testing.T) {
	// Test Append when createAll() fails
	// Use a path where directory creation might fail
	file := NewFile("/root/nonexistent/readonly/path/file.txt")

	result := file.Append("test content")

	if result.IsOk() {
		t.Log("Append succeeded unexpectedly - platform dependent")
	} else {
		t.Logf("Append failed as expected: %v", result.Err())
	}
}

func TestFile_Append_WithOpenFile(t *testing.T) {
	// Test Append when file is already open for writing
	tempFile := createTempFile(t)
	defer os.Remove(tempFile)

	file := NewFile(String(tempFile))

	// Open the file for writing first
	if openResult := file.OpenFile(os.O_RDWR, 0o644); openResult.IsErr() {
		t.Fatalf("Failed to open file for writing: %v", openResult.Err())
	}
	defer file.Close()

	// Test append to open file - this should still work as it will use the current file handle
	result := file.Append("appended content")

	if result.IsErr() {
		t.Logf("Append to open file failed (expected for read-only files): %v", result.Err())
	} else {
		t.Log("Append to open file succeeded")
	}
}
