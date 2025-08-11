package g_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/enetx/g"
)

func TestNewDir(t *testing.T) {
	// Define a test path
	testPath := String("/")

	// Create a new Dir instance using NewDir
	dir := NewDir(testPath)

	// Check if the path of the created Dir instance matches the expected path
	if dir.Path().Ok() != testPath {
		t.Errorf("TestNewDir: Expected path %s, got %s", testPath, dir.Path().Ok())
	}
}

func TestDir_Chown_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Perform chown operation
	uid := os.Getuid() // Current user's UID
	gid := os.Getgid() // Current user's GID
	result := dir.Chown(uid, gid)

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Chown_Success: Unexpected error: %s", result.Err().Error())
	}
}

func TestDir_Chown_Failure(t *testing.T) {
	// Create a Dir instance with a non-existent directory path
	dir := NewDir("/nonexistent/path")

	// Perform chown operation
	uid := os.Getuid() // Current user's UID
	gid := os.Getgid() // Current user's GID
	result := dir.Chown(uid, gid)

	// Check if the operation failed as expected
	if result.IsOk() {
		t.Errorf("TestDir_Chown_Failure: Expected error, got success")
	}
}

func TestDir_Stat_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Get directory information using Stat
	result := dir.Stat()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Stat_Success: Unexpected error: %s", result.Err().Error())
	}
}

func TestDir_Stat_Failure(t *testing.T) {
	// Create a Dir instance with a non-existent directory path
	dir := NewDir("/nonexistent/path")

	// Get directory information using Stat
	result := dir.Stat()

	// Check if the operation failed as expected
	if result.IsOk() {
		t.Errorf("TestDir_Stat_Failure: Expected error, got success")
	}
}

func TestDir_Path_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Get the absolute path of the directory using Path
	result := dir.Path()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Path_Success: Unexpected error: %s", result.Err().Error())
	}
}

func TestDir_Lstat_IsLink_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a symbolic link to the temporary directory
	linkPath := tempDir + "/link"
	err := os.Symlink(tempDir, linkPath)
	if err != nil {
		t.Fatalf("Failed to create symbolic link: %s", err)
	}

	// Create a Dir instance with the symbolic link path
	dir := NewDir(String(linkPath))

	// Call Lstat to get information about the symbolic link
	result := dir.Lstat()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Lstat_IsLink_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if IsLink correctly identifies the symbolic link
	if !dir.IsLink() {
		t.Errorf("TestDir_Lstat_IsLink_Success: Expected directory to be a symbolic link")
	}
}

func TestDir_Lstat_IsLink_NotLink(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Call Lstat to get information about the directory
	result := dir.Lstat()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Lstat_IsLink_NotLink: Unexpected error: %s", result.Err().Error())
	}

	// Check if IsLink correctly identifies that it's not a symbolic link
	if dir.IsLink() {
		t.Errorf("TestDir_Lstat_IsLink_NotLink: Expected directory not to be a symbolic link")
	}
}

func TestDir_CreateTemp_Success(t *testing.T) {
	// Create a Dir instance representing the default directory for temporary directories
	dir := NewDir("")

	// Create a temporary directory using CreateTemp
	result := dir.CreateTemp()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_CreateTemp_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the temporary directory exists
	tmpDir := result.Ok().Path().Ok().Std()
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Errorf("TestDir_CreateTemp_Success: Temporary directory not created")
	}
}

func TestDir_Temp(t *testing.T) {
	// Get the default temporary directory using Temp
	tmpDir := NewDir("").Temp()

	// Check if the returned directory exists
	if _, err := os.Stat(tmpDir.String().Std()); os.IsNotExist(err) {
		t.Errorf("TestDir_Temp: Temporary directory does not exist")
	}
}

func TestDir_Remove_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Remove the temporary directory using Remove
	result := dir.Remove()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Remove_Success: Unexpected error: %s", result.Err().Error())
	}
}

func TestDir_Remove_NotExist(t *testing.T) {
	// Create a Dir instance with a non-existent directory path
	dir := NewDir("/nonexistent/path")

	// Remove the non-existent directory using Remove
	result := dir.Remove()

	// Check if the operation succeeded (non-existent directory should be considered removed)
	if result.IsErr() {
		t.Errorf("TestDir_Remove_NotExist: Unexpected error: %s", result.Err().Error())
	}
}

func TestDir_Create_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	os.RemoveAll(tempDir)

	defer os.RemoveAll(tempDir)

	// Create a Dir instance with the temporary directory path
	dir := NewDir(String(tempDir))

	// Create a new directory using Create
	result := dir.Create()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Create_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the created directory exists
	createdDir := dir.Path().Ok().Std()
	if _, err := os.Stat(createdDir); os.IsNotExist(err) {
		t.Errorf("TestDir_Create_Success: Created directory does not exist")
	}
}

func TestDir_Create_Failure(t *testing.T) {
	// Attempt to create a directory in a non-existent parent directory
	nonExistentDir := NewDir("/nonexistent/parent")
	result := nonExistentDir.Create()

	// Check if the operation failed as expected
	if result.IsOk() {
		t.Errorf("TestDir_Create_Failure: Expected error, got success")
	}
}

func TestDir_Join_Success(t *testing.T) {
	// Create a Dir instance representing an existing directory
	dir := NewDir("/path/to/directory")

	// Join the directory path with additional elements
	result := dir.Join("subdir", "file.txt")

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Join_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the joined path matches the expected value
	expectedPath := "/path/to/directory/subdir/file.txt"
	if result.Ok().Std() != expectedPath {
		t.Errorf("TestDir_Join_Success: Expected joined path '%s', got '%s'", expectedPath, result.Ok().Std())
	}
}

func TestDir_SetPath(t *testing.T) {
	// Create a Dir instance representing an existing directory
	dir := NewDir("/path/to/directory")

	// Set a new path for the directory
	newPath := String("/new/path/to/directory")
	updatedDir := dir.SetPath(newPath)

	// Check if the path of the directory is updated correctly
	if updatedDir.Path().Ok() != newPath {
		t.Errorf("TestDir_SetPath: Expected path '%s', got '%s'", newPath, updatedDir.Path().Ok())
	}
}

func TestDir_CreateAll_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance representing the temporary directory
	dir := NewDir(String(tempDir))

	// Create all directories along the path
	result := dir.CreateAll()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_CreateAll_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the directories along the path are created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("TestDir_CreateAll_Success: Directory does not exist: %s", tempDir)
	}
}

func TestDir_CreateAll_Mode_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance representing the temporary directory
	dir := NewDir(String(tempDir))

	// Create all directories along the path with custom mode
	result := dir.CreateAll(0700)

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_CreateAll_Mode_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the directories along the path are created with the specified mode
	fileInfo, err := os.Stat(tempDir)
	if os.IsNotExist(err) {
		t.Errorf("TestDir_CreateAll_Mode_Success: Directory does not exist: %s", tempDir)
	} else {
		if fileInfo.Mode() != os.FileMode(0700)|os.ModeDir {
			t.Errorf("TestDir_CreateAll_Mode_Success: Expected mode 0700, got %o", fileInfo.Mode())
		}
	}
}

func TestDir_Read_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create some files and directories inside the temporary directory
	createTestFiles(tempDir)

	// Create a Dir instance representing the temporary directory
	dir := NewDir(String(tempDir))

	// Read the content of the directory
	result := dir.Read().Collect()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Read_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the returned slice of File instances is accurate
	files := result.Ok()
	expectedFileNames := []string{"file1.txt", "file2.txt", "subdir1", "subdir2"}

	for i, file := range files {
		if file.Name().Std() != expectedFileNames[i] {
			t.Errorf("TestDir_Read_Success: Expected file '%s', got '%s'", expectedFileNames[i], file.Name().Std())
		}
	}
}

func TestDir_Glob_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create some test files inside the temporary directory
	createTestFiles(tempDir)

	// Create a Dir instance representing the temporary directory with a glob pattern
	dir := NewDir(String(filepath.Join(tempDir, "*.txt")))

	// Retrieve files matching the glob pattern
	result := dir.Glob().Collect()

	// Check if the operation succeeded
	if result.IsErr() {
		t.Errorf("TestDir_Glob_Success: Unexpected error: %s", result.Err().Error())
	}

	// Check if the returned slice of File instances is accurate
	files := result.Ok()
	expectedFileNames := []string{"file1.txt", "file2.txt"}
	for i, file := range files {
		if file.Name().Std() != expectedFileNames[i] {
			t.Errorf("TestDir_Glob_Success: Expected file '%s', got '%s'", expectedFileNames[i], file.Name().Std())
		}
	}
}

func TestDir_Rename(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	// Create a Dir instance representing the temporary directory
	dir := NewDir(String(tempDir))

	// Rename the directory
	newpath := tempDir + "_new"
	renamedDir := dir.Rename(String(newpath))

	// Check if the directory has been successfully renamed
	if renamedDir.IsErr() {
		t.Errorf("TestDir_Rename: Failed to rename directory: %v", renamedDir.Err())
	}

	defer renamedDir.Ok().Remove()

	// Verify that the old directory does not exist
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("TestDir_Rename: Old directory still exists after renaming")
	}

	// Verify that the new directory exists
	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		t.Errorf("TestDir_Rename: New directory does not exist after renaming")
	}
}

func TestDir_Copy(t *testing.T) {
	// Create a temporary source directory for testing
	sourceDir := createTempDir(t)
	defer os.RemoveAll(sourceDir)

	// Create some test files in the source directory
	if err := os.WriteFile(sourceDir+"/file1.txt", []byte("File 1 content"), 0644); err != nil {
		t.Fatalf("TestDir_Copy: Failed to create test file 1 in source directory: %v", err)
	}
	if err := os.WriteFile(sourceDir+"/file2.txt", []byte("File 2 content"), 0644); err != nil {
		t.Fatalf("TestDir_Copy: Failed to create test file 2 in source directory: %v", err)
	}

	// Create a temporary destination directory for testing
	destDir := createTempDir(t)
	defer os.RemoveAll(destDir)

	// Create a Dir instance representing the source directory
	source := NewDir(String(sourceDir))

	// Copy the contents of the source directory to the destination directory
	result := source.Copy(String(destDir))
	if result.IsErr() {
		t.Fatalf("TestDir_Copy: Failed to copy directory contents: %v", result.Err())
	}

	// Verify that the destination directory contains the same files as the source directory
	destFiles, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatalf("TestDir_Copy: Failed to read destination directory: %v", err)
	}

	expectedFiles := []string{"file1.txt", "file2.txt"}
	for _, expectedFile := range expectedFiles {
		found := false
		for _, destFile := range destFiles {
			if destFile.Name() == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("TestDir_Copy: Destination directory missing file %s", expectedFile)
		}
	}
}

func TestDir_Walk(t *testing.T) {
	testDir := createTempDir(t)
	defer os.RemoveAll(testDir)

	createTestFiles(testDir)

	var visited Slice[string]

	NewDir(String(testDir)).Walk().
		Exclude(func(f *File) bool { return f.IsDir() && f.Dir().Ok().IsLink() }).
		ForEach(func(r Result[*File]) {
			if r.IsErr() {
				t.Fatalf("unexpected error during Walk: %v", r.Err())
			}
			visited.Push(r.Ok().Path().Ok().Std())
		})

	expected := []string{
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "subdir1"),
		filepath.Join(testDir, "subdir2"),
	}

	for _, exp := range expected {
		if !visited.Contains(exp) {
			t.Errorf("TestDir_Walk: expected %q to be visited, got %v", exp, visited)
		}
	}
}

func createTestFiles(dir string) {
	os.Mkdir(filepath.Join(dir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(dir, "subdir2"), 0755)
	os.Create(filepath.Join(dir, "file1.txt"))
	os.Create(filepath.Join(dir, "file2.txt"))
}

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %s", err)
	}
	return tempDir
}

func TestDirMove(t *testing.T) {
	// Create source directory
	srcDir := NewDir(String(createTempDir(t)))
	defer srcDir.Remove()

	// Create a test file in source directory
	testFilePath := srcDir.Join("testfile.txt")
	if testFilePath.IsErr() {
		t.Fatalf("Failed to join path: %v", testFilePath.Err())
	}
	testFile := NewFile(testFilePath.Ok())
	testFile.Write(String("test content"))

	// Create destination directory path
	destPath := String(createTempDir(t) + "_moved")
	destDir := NewDir(destPath)
	defer destDir.Remove()

	// Move directory
	result := srcDir.Move(destPath)
	if result.IsErr() {
		t.Errorf("Move should succeed, got error: %v", result.Err())
	}

	// Check that destination exists
	if !destDir.Exist() {
		t.Error("Destination directory should exist after move")
	}

	// Check that source doesn't exist
	if srcDir.Exist() {
		t.Error("Source directory should not exist after move")
	}

	// Check that file was moved
	movedFilePath := destDir.Join("testfile.txt")
	if movedFilePath.IsErr() {
		t.Fatalf("Failed to join moved file path: %v", movedFilePath.Err())
	}
	movedFile := NewFile(movedFilePath.Ok())
	if !movedFile.Exist() {
		t.Error("File should exist in moved directory")
	}
}

func TestDirPrint(t *testing.T) {
	dir := NewDir("testdir")
	result := dir.Print()
	if result != dir {
		t.Errorf("Print() should return original dir unchanged")
	}
}

func TestDirPrintln(t *testing.T) {
	dir := NewDir("testdir")
	result := dir.Println()
	if result != dir {
		t.Errorf("Println() should return original dir unchanged")
	}
}
