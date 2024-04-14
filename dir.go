package g

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// SkipWalk is used as a return value from the walker function to indicate that
// the file or directory named in the call should be skipped. It is not returned
// as an error by any function.
var SkipWalk = errors.New("skip")

// StopWalk is used as a return value from the walker function to indicate that
// all remaining files and directories should be skipped. It is not returned
// as an error by any function.
var StopWalk = errors.New("stop")

// NewDir returns a new Dir instance with the given path.
func NewDir(path String) *Dir { return &Dir{path: path} }

// Chown changes the ownership of the directory to the specified UID and GID.
// It uses os.Chown to modify ownership and returns a Result[*Dir] indicating success or failure.
func (d *Dir) Chown(uid, gid int) Result[*Dir] {
	err := os.Chown(d.path.Std(), uid, gid)
	if err != nil {
		return Err[*Dir](err)
	}

	return Ok(d)
}

// Stat retrieves information about the directory represented by the Dir instance.
// It returns a Result[fs.FileInfo] containing details about the directory's metadata.
func (d *Dir) Stat() Result[fs.FileInfo] {
	if d.Path().IsErr() {
		return Err[fs.FileInfo](d.Path().Err())
	}

	return ResultOf(os.Stat(d.Path().Ok().Std()))
}

// Lstat retrieves information about the symbolic link represented by the Dir instance.
// It returns a Result[fs.FileInfo] containing details about the symbolic link's metadata.
// Unlike Stat, Lstat does not follow the link and provides information about the link itself.
func (d *Dir) Lstat() Result[fs.FileInfo] {
	return ResultOf(os.Lstat(d.Path().Ok().Std()))
}

// IsLink checks if the directory is a symbolic link.
func (d *Dir) IsLink() bool {
	stat := d.Lstat()
	return stat.IsOk() && stat.Ok().Mode()&os.ModeSymlink != 0
}

// CreateTemp creates a new temporary directory in the specified directory with the
// specified name pattern and returns a Result, which contains a pointer to the Dir
// or an error if the operation fails.
// If no directory is specified, the default directory for temporary directories is used.
// If no name pattern is specified, the default pattern "*" is used.
//
// Parameters:
//
// - args ...String: A variadic parameter specifying the directory and/or name
// pattern for the temporary directory.
//
// Returns:
//
// - *Dir: A pointer to the Dir representing the temporary directory.
//
// Example usage:
//
// d := g.NewDir("")
// tmpdir := d.CreateTemp()                     // Creates a temporary directory with default settings
// tmpdirWithDir := d.CreateTemp("mydir")       // Creates a temporary directory in "mydir" directory
// tmpdirWithPattern := d.CreateTemp("", "tmp") // Creates a temporary directory with "tmp" pattern
func (*Dir) CreateTemp(args ...String) Result[*Dir] {
	dir := ""
	pattern := "*"

	if len(args) != 0 {
		if len(args) > 1 {
			pattern = args[1].Std()
		}

		dir = args[0].Std()
	}

	tmpDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return Err[*Dir](err)
	}

	return Ok(NewDir(String(tmpDir)))
}

// Temp returns the default directory to use for temporary files.
//
// On Unix systems, it returns $TMPDIR if non-empty, else /tmp.
// On Windows, it uses GetTempPath, returning the first non-empty
// value from %TMP%, %TEMP%, %USERPROFILE%, or the Windows directory.
// On Plan 9, it returns /tmp.
//
// The directory is neither guaranteed to exist nor have accessible
// permissions.
func (*Dir) Temp() *Dir { return NewDir(String(os.TempDir())) }

// Remove attempts to delete the directory and its contents.
// It returns a Result, which contains either the *Dir or an error.
// If the directory does not exist, Remove returns a successful Result with *Dir set.
// Any error that occurs during removal will be of type *PathError.
func (d *Dir) Remove() Result[*Dir] {
	if err := os.RemoveAll(d.ToString().Std()); err != nil {
		return Err[*Dir](err)
	}

	return Ok(d)
}

// Copy copies the contents of the current directory to the destination directory.
//
// Parameters:
//
// - dest (String): The destination directory where the contents of the current directory should be copied.
//
// - followLinks (optional): A boolean indicating whether to follow symbolic links during the walk.
// If true, symbolic links are followed; otherwise, they are skipped.
//
// Returns:
//
// - Result[*Dir]: A Result type containing either a pointer to a new Dir instance representing the destination directory or an error.
//
// Example usage:
//
//	sourceDir := g.NewDir("path/to/source")
//	destinationDirResult := sourceDir.Copy("path/to/destination")
//	if destinationDirResult.IsErr() {
//		// Handle error
//	}
//	destinationDir := destinationDirResult.Ok()
func (d *Dir) Copy(dest String, followLinks ...bool) Result[*Dir] {
	follow := true
	if len(followLinks) != 0 {
		follow = followLinks[0]
	}

	root := d.Path()
	if root.IsErr() {
		return Err[*Dir](root.Err())
	}

	walker := func(f *File) error {
		path := f.Path()
		if path.IsErr() {
			return path.Err()
		}

		relpath, err := filepath.Rel(root.Ok().Std(), path.Ok().Std())
		if err != nil {
			return err
		}

		destpath := NewDir(dest).Join(String(relpath))
		if destpath.IsErr() {
			return destpath.Err()
		}

		stat := f.Stat()
		if stat.IsErr() {
			return stat.Err()
		}

		if stat.Ok().IsDir() {
			if !follow && f.IsLink() {
				return SkipWalk
			}

			return NewDir(destpath.Ok()).CreateAll(stat.Ok().Mode()).Err()
		}

		return f.Copy(destpath.Ok(), stat.Ok().Mode()).Err()
	}

	if err := d.Walk(walker); err != nil {
		return Err[*Dir](err)
	}

	return Ok(NewDir(dest))
}

// Create creates a new directory with the specified mode (optional).
//
// Parameters:
//
// - mode (os.FileMode, optional): The file mode for the new directory.
// If not provided, it defaults to DirDefault (0755).
//
// Returns:
//
// - *Dir: A pointer to the Dir instance on which the method was called.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	createdDir := dir.Create(0755) // Optional mode argument
func (d *Dir) Create(mode ...os.FileMode) Result[*Dir] {
	dmode := DirDefault
	if len(mode) != 0 {
		dmode = mode[0]
	}

	if err := os.Mkdir(d.path.Std(), dmode); err != nil {
		return Err[*Dir](err)
	}

	return Ok(d)
}

// Join joins the current directory path with the given path elements, returning the joined path.
//
// Parameters:
//
// - elem (...String): One or more String values representing path elements to
// be joined with the current directory path.
//
// Returns:
//
// - String: The resulting joined path as an String.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	joinedPath := dir.Join("subdir", "file.txt")
func (d *Dir) Join(elem ...String) Result[String] {
	path := d.Path()
	if path.IsErr() {
		return Err[String](path.Err())
	}

	paths := SliceOf(elem...).Insert(0, path.Ok()).ToStringSlice()

	return Ok(String(filepath.Join(paths...)))
}

// SetPath sets the path of the current directory.
//
// Parameters:
//
// - path (String): The new path to be set for the current directory.
//
// Returns:
//
// - *Dir: A pointer to the updated Dir instance with the new path.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	dir.SetPath("new/path/to/directory")
func (d *Dir) SetPath(path String) *Dir {
	d.path = path
	return d
}

// CreateAll creates all directories along the given path, with the specified mode (optional).
//
// Parameters:
//
// - mode ...os.FileMode (optional): The file mode to be used when creating the directories.
// If not provided, it defaults to the value of DirDefault constant (0755).
//
// Returns:
//
// - *Dir: A pointer to the Dir instance representing the created directories.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	dir.CreateAll()
//	dir.CreateAll(0755)
func (d *Dir) CreateAll(mode ...os.FileMode) Result[*Dir] {
	if d.Exist() {
		return Ok(d)
	}

	dmode := DirDefault
	if len(mode) != 0 {
		dmode = mode[0]
	}

	path := d.Path()
	if path.IsErr() {
		return Err[*Dir](path.Err())
	}

	err := os.MkdirAll(path.Ok().Std(), dmode)
	if err != nil {
		return Err[*Dir](err)
	}

	return Ok(d)
}

// Rename renames the current directory to the new path.
//
// Parameters:
//
// - newpath String: The new path for the directory.
//
// Returns:
//
// - *Dir: A pointer to the Dir instance representing the renamed directory.
// If an error occurs, the original Dir instance is returned with the error stored in d.err,
// which can be checked using the Error() method.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	dir.Rename("path/to/new_directory")
func (d *Dir) Rename(newpath String) Result[*Dir] {
	ps := String(os.PathSeparator)
	_, np := newpath.TrimSuffix(ps).Split(ps).Pop()

	if rd := NewDir(np.Join(ps)).CreateAll(); rd.IsErr() {
		return rd
	}

	if err := os.Rename(d.path.Std(), newpath.Std()); err != nil {
		return Err[*Dir](err)
	}

	return Ok(NewDir(newpath))
}

// Move function simply calls [Dir.Rename]
func (d *Dir) Move(newpath String) Result[*Dir] { return d.Rename(newpath) }

// Path returns the absolute path of the current directory.
//
// Returns:
//
// - String: The absolute path of the current directory as an String.
// If an error occurs while converting the path to an absolute path,
// the error is stored in d.err, which can be checked using the Error() method.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	absPath := dir.Path()
func (d *Dir) Path() Result[String] {
	path, err := filepath.Abs(d.path.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(path))
}

// Exist checks if the current directory exists.
//
// Returns:
//
// - bool: true if the current directory exists, false otherwise.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	exists := dir.Exist()
func (d *Dir) Exist() bool {
	path := d.Path()
	if path.IsErr() {
		return false
	}

	_, err := os.Stat(path.Ok().Std())

	return !os.IsNotExist(err)
}

// Read reads the content of the current directory and returns a slice of File instances.
//
// Returns:
//
// - []*File: A slice of File instances representing the files and directories
// in the current directory.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory")
//	files := dir.Read()

func (d *Dir) Read() Result[Slice[*File]] {
	entries, err := os.ReadDir(d.path.Std())
	if err != nil {
		return Err[Slice[*File]](err)
	}

	files := NewSlice[*File](0, Int(len(entries)))

	for _, entry := range entries {
		dpath := d.Path()
		if dpath.IsErr() {
			return Err[Slice[*File]](dpath.Err())
		}

		file := NewDir(dpath.Ok()).Join(String(entry.Name()))
		if file.IsErr() {
			return Err[Slice[*File]](file.Err())
		}

		files = files.Append(NewFile(file.Ok()))
	}

	return Ok(files)
}

// Glob matches files in the current directory using the path pattern and
// returns a slice of File instances.
//
// Returns:
//
// - []*File: A slice of File instances representing the files that match the
// provided pattern in the current directory.
//
// Example usage:
//
//	dir := g.NewDir("path/to/directory/*.txt")
//	files := dir.Glob()
func (d *Dir) Glob() Result[Slice[*File]] {
	matches, err := filepath.Glob(d.path.Std())
	if err != nil {
		return Err[Slice[*File]](err)
	}

	files := NewSlice[*File](0, Int(len(matches)))

	for _, match := range matches {
		file := NewFile(String(match)).Path()
		if file.IsErr() {
			return Err[Slice[*File]](file.Err())
		}

		files = files.Append(NewFile(file.Ok()))
	}

	return Ok(files)
}

// Walk recursively traverses the directory structure and applies the walker function to each file and directory.
//
// Parameters:
//
// - walker: A function that takes a *File as an argument and returns an error.
// It is applied to each file and directory encountered during the walk.
//
// Returns:
//
// - error: An error indicating any issues that occurred during the walk. If no errors occurred, it returns nil.
func (d *Dir) Walk(walker func(f *File) error) error {
	entries := d.Read()
	if entries.IsErr() {
		return entries.Err()
	}

	for _, entry := range entries.Ok() {
		if err := walker(entry); err != nil {
			switch {
			case errors.Is(err, SkipWalk):
				continue
			case errors.Is(err, StopWalk):
				return nil
			default:
				return err
			}
		}

		stat := entry.Stat()
		if stat.IsErr() {
			return stat.Err()
		}

		if stat.Ok().IsDir() {
			entryPath := entry.Path()
			if entryPath.IsErr() {
				return entryPath.Err()
			}

			subdir := NewDir(entryPath.Ok())

			if err := subdir.Walk(walker); err != nil {
				return err
			}
		}
	}

	return nil
}

// ToString returns the String representation of the current directory's path.
func (d *Dir) ToString() String { return d.path }

// Print prints the content of the Dir to the standard output (console)
// and returns the Dir unchanged.
func (d *Dir) Print() *Dir { fmt.Println(d); return d }
