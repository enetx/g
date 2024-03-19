package g_test

import (
	"os"
	"path/filepath"
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

// createTempFile creates a temporary file for testing and returns its path.
func createTempFile(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	return tempFile.Name()
}
