package g

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/enetx/g/internal/filelock"
)

// NewFile returns a new File instance with the given name.
func NewFile(name String) *File { return &File{name: name} }

// Lines returns a new iterator instance that can be used to read the file
// line by line.
//
// Example usage:
//
//	// Open a new file with the specified name "text.txt"
//	g.NewFile("text.txt").
//		Lines().                 // Read the file line by line
//		Unwrap().                // Unwrap the Result type to get the underlying iterator
//		Skip(3).                 // Skip the first 3 lines
//		Exclude(filters.IsZero). // Exclude lines that are empty or contain only whitespaces
//		Dedup().                 // Remove consecutive duplicate lines
//		Map(g.String.Upper).     // Convert each line to uppercase
//		ForEach(                 // For each line, print it
//			func(s g.String) {
//				s.Print()
//			})
//
//	// Output:
//	// UPPERCASED_LINE4
//	// UPPERCASED_LINE5
//	// UPPERCASED_LINE6
func (f *File) Lines() Result[seqSlice[String]] {
	if f.file == nil {
		if err := f.Open().Err(); err != nil {
			return Err[seqSlice[String]](err)
		}
	}

	return Ok(seqSlice[String](func(yield func(String) bool) {
		defer f.Close()

		scanner := bufio.NewScanner(f.file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			if !yield(String(scanner.Text())) {
				return
			}
		}
	}))
}

// Chunks returns a new iterator instance that can be used to read the file
// in fixed-size chunks of the specified size in bytes.
//
// Parameters:
//
// - size (int): The size of each chunk in bytes.
//
// Example usage:
//
//	// Open a new file with the specified name "text.txt"
//	g.NewFile("text.txt").
//		Chunks(100).              // Read the file in chunks of 100 bytes
//		Unwrap().                 // Unwrap the Result type to get the underlying iterator
//		Map(g.String.Upper).      // Convert each chunk to uppercase
//		ForEach(                  // For each chunk, print it
//			func(s g.String) {
//				s.Print()
//			})
//
//	// Output:
//	// UPPERCASED_CHUNK1
//	// UPPERCASED_CHUNK2
//	// UPPERCASED_CHUNK3
func (f *File) Chunks(size int) Result[seqSlice[String]] {
	if f.file == nil {
		if err := f.Open().Err(); err != nil {
			return Err[seqSlice[String]](err)
		}
	}

	return Ok(seqSlice[String](func(yield func(String) bool) {
		defer f.Close()

		buffer := make([]byte, size)
		for {
			n, err := f.file.Read(buffer)
			if err != nil && err != io.EOF {
				return
			}

			if n == 0 {
				break
			}

			if !yield(String(buffer[:n])) {
				return
			}
		}
	}))
}

// Append appends the given content to the file, with the specified mode (optional).
// If no FileMode is provided, the default FileMode (0644) is used.
// Don't forget to close the file!
func (f *File) Append(content String, mode ...os.FileMode) Result[*File] {
	if f.file == nil {
		if mda := f.createAll(); mda.IsErr() {
			return mda
		}

		fmode := FileDefault
		if len(mode) != 0 {
			fmode = mode[0]
		}

		if err := f.OpenFile(os.O_APPEND|os.O_CREATE|os.O_WRONLY, fmode).Err(); err != nil {
			return Err[*File](err)
		}
	}

	if _, err := f.file.WriteString(content.Std()); err != nil {
		return Err[*File](err)
	}

	return Ok(f)
}

// Chmod changes the mode of the file.
func (f *File) Chmod(mode os.FileMode) Result[*File] {
	var err error
	if f.file != nil {
		err = f.file.Chmod(mode)
	} else {
		err = os.Chmod(f.name.Std(), mode)
	}

	if err != nil {
		return Err[*File](err)
	}

	return Ok(f)
}

// Chown changes the owner of the file.
func (f *File) Chown(uid, gid int) Result[*File] {
	var err error
	if f.file != nil {
		err = f.file.Chown(uid, gid)
	} else {
		err = os.Chown(f.name.Std(), uid, gid)
	}

	if err != nil {
		return Err[*File](err)
	}

	return Ok(f)
}

// Seek sets the file offset for the next Read or Write operation. The offset
// is specified by the 'offset' parameter, and the 'whence' parameter determines
// the reference point for the offset.
//
// The 'offset' parameter specifies the new offset in bytes relative to the
// reference point determined by 'whence'. If 'whence' is set to io.SeekStart,
// io.SeekCurrent, or io.SeekEnd, the offset is relative to the start of the file,
// the current offset, or the end of the file, respectively.
//
// If the file is not open, this method will attempt to open it. If the open
// operation fails, an error is returned.
//
// If the Seek operation fails, the file is closed, and an error is returned.
//
// Example:
//
//	file := g.NewFile("example.txt")
//	result := file.Seek(100, io.SeekStart)
//	if result.Err() != nil {
//	    log.Fatal(result.Err())
//	}
//
// Parameters:
//   - offset: The new offset in bytes.
//   - whence: The reference point for the offset (io.SeekStart, io.SeekCurrent, or io.SeekEnd).
//
// Don't forget to close the file!
func (f *File) Seek(offset int64, whence int) Result[*File] {
	if f.file == nil {
		if err := f.Open().Err(); err != nil {
			return Err[*File](err)
		}
	}

	if _, err := f.file.Seek(offset, whence); err != nil {
		f.Close()
		return Err[*File](err)
	}

	return Ok(f)
}

// Close closes the File and unlocks its underlying file, if it is not already closed.
func (f *File) Close() error {
	if f.file == nil {
		return &ErrFileClosed{f.name.Std()}
	}

	var err error

	if f.guard {
		err = filelock.Unlock(f.file)
	}

	if closeErr := f.file.Close(); closeErr != nil {
		err = closeErr
	}

	f.file = nil

	return err
}

// Copy copies the file to the specified destination, with the specified mode (optional).
// If no mode is provided, the default FileMode (0644) is used.
func (f *File) Copy(dest String, mode ...os.FileMode) Result[*File] {
	if err := f.Open().Err(); err != nil {
		return Err[*File](err)
	}
	defer f.Close()

	return NewFile(dest).WriteFromReader(f.file, mode...)
}

// Create is similar to os.Create, but it returns a write-locked file.
// Don't forget to close the file!
func (f *File) Create() Result[*File] {
	return f.OpenFile(os.O_RDWR|os.O_CREATE|os.O_TRUNC, FileCreate)
}

// Dir returns the directory the file is in as an Dir instance.
func (f *File) Dir() Result[*Dir] {
	dirPath := f.dirPath()
	if dirPath.IsErr() {
		return Err[*Dir](dirPath.Err())
	}

	return Ok(NewDir(dirPath.Ok()))
}

// Exist checks if the file exists.
func (f *File) Exist() bool {
	if f.dirPath().IsOk() {
		filePath := f.filePath()
		if filePath.IsOk() {
			_, err := os.Stat(filePath.Ok().Std())
			return !os.IsNotExist(err)
		}
	}

	return false
}

// Ext returns the file extension.
func (f *File) Ext() String { return String(filepath.Ext(f.name.Std())) }

// Guard sets a lock on the file to protect it from concurrent access.
// It returns the File instance with the guard enabled.
func (f *File) Guard() *File {
	f.guard = true
	return f
}

// MimeType returns the MIME type of the file as an String.
func (f *File) MimeType() Result[String] {
	if err := f.Open().Err(); err != nil {
		return Err[String](err)
	}
	defer f.Close()

	const bufferSize = 512

	buff := make([]byte, bufferSize)

	bytesRead, err := f.file.ReadAt(buff, 0)
	if err != nil && err != io.EOF {
		return Err[String](err)
	}

	buff = buff[:bytesRead]

	return Ok(String(http.DetectContentType(buff)))
}

// Move function simply calls [File.Rename]
func (f *File) Move(newpath String) Result[*File] { return f.Rename(newpath) }

// Name returns the name of the file.
func (f *File) Name() String {
	if f.file != nil {
		return String(filepath.Base(f.file.Name()))
	}

	return String(filepath.Base(f.name.Std()))
}

// Open is like os.Open, but returns a read-locked file.
// Don't forget to close the file!
func (f *File) Open() Result[*File] { return f.OpenFile(os.O_RDONLY, 0) }

// OpenFile is like os.OpenFile, but returns a locked file.
// If flag includes os.O_WRONLY or os.O_RDWR, the file is write-locked
// otherwise, it is read-locked.
// Don't forget to close the file!
func (f *File) OpenFile(flag int, perm fs.FileMode) Result[*File] {
	file, err := os.OpenFile(f.name.Std(), flag&^os.O_TRUNC, perm)
	if err != nil {
		return Err[*File](err)
	}

	if f.guard {
		switch flag & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR) {
		case os.O_WRONLY, os.O_RDWR:
			err = filelock.Lock(file)
		default:
			err = filelock.RLock(file)
		}

		if err != nil {
			file.Close()
			return Err[*File](err)
		}
	}

	if flag&os.O_TRUNC == os.O_TRUNC {
		if err := file.Truncate(0); err != nil {
			if fi, statErr := file.Stat(); statErr != nil || fi.Mode().IsRegular() {
				if f.guard {
					filelock.Unlock(file)
				}

				file.Close()

				return Err[*File](err)
			}
		}
	}

	f.file = file

	return Ok(f)
}

// Path returns the absolute path of the file.
func (f *File) Path() Result[String] { return f.filePath() }

// Print prints the content of the File to the standard output (console)
// and returns the File unchanged.
func (f *File) Print() *File { fmt.Println(f); return f }

// Read opens the named file with a read-lock and returns its contents.
func (f *File) Read() Result[String] {
	if err := f.Open().Err(); err != nil {
		return Err[String](err)
	}

	defer f.Close()

	content, err := io.ReadAll(f.file)
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(content))
}

// Remove removes the file.
func (f *File) Remove() Result[*File] {
	if err := os.Remove(f.name.Std()); err != nil {
		return Err[*File](err)
	}

	return Ok(f)
}

// Rename renames the file to the specified new path.
func (f *File) Rename(newpath String) Result[*File] {
	if !f.Exist() {
		return Err[*File](&ErrFileNotExist{f.name.Std()})
	}

	nf := NewFile(newpath).createAll()

	if err := os.Rename(f.name.Std(), newpath.Std()); err != nil {
		return Err[*File](err)
	}

	return nf
}

// Split splits the file path into its directory and file components.
func (f *File) Split() (*Dir, *File) {
	path := f.Path()
	if path.IsErr() {
		return nil, nil
	}

	dir, file := filepath.Split(path.Ok().Std())

	return NewDir(String(dir)), NewFile(String(file))
}

// Stat returns the fs.FileInfo of the file.
// It calls the file's Stat method if the file is open, or os.Stat otherwise.
func (f *File) Stat() Result[fs.FileInfo] {
	if f.file != nil {
		return ResultOf(f.file.Stat())
	}

	return ResultOf(os.Stat(f.name.Std()))
}

// Lstat retrieves information about the symbolic link represented by the *File instance.
// It returns a Result[fs.FileInfo] containing details about the symbolic link's metadata.
// Unlike Stat, Lstat does not follow the link and provides information about the link itself.
func (f *File) Lstat() Result[fs.FileInfo] {
	return ResultOf(os.Lstat(f.name.Std()))
}

// IsDir checks if the file is a directory.
func (f *File) IsDir() bool {
	stat := f.Stat()
	return stat.IsOk() && stat.Ok().IsDir()
}

// IsLink checks if the file is a symbolic link.
func (f *File) IsLink() bool {
	stat := f.Lstat()
	return stat.IsOk() && stat.Ok().Mode()&os.ModeSymlink != 0
}

// Std returns the underlying *os.File instance.
// Don't forget to close the file with g.File().Close()!
func (f *File) Std() *os.File { return f.file }

// CreateTemp creates a new temporary file in the specified directory with the
// specified name pattern and returns a Result, which contains a pointer to the File
// or an error if the operation fails.
// If no directory is specified, the default directory for temporary files is used.
// If no name pattern is specified, the default pattern "*" is used.
//
// Parameters:
//
// - args ...String: A variadic parameter specifying the directory and/or name
// pattern for the temporary file.
//
// Returns:
//
// - *File: A pointer to the File representing the temporary file.
//
// Example usage:
//
//	f := g.NewFile("")
//	tmpfile := f.CreateTemp()                     // Creates a temporary file with default settings
//	tmpfileWithDir := f.CreateTemp("mydir")       // Creates a temporary file in "mydir" directory
//	tmpfileWithPattern := f.CreateTemp("", "tmp") // Creates a temporary file with "tmp" pattern
func (*File) CreateTemp(args ...String) Result[*File] {
	dir := ""
	pattern := "*"

	if len(args) != 0 {
		if len(args) > 1 {
			pattern = args[1].Std()
		}

		dir = args[0].Std()
	}

	tmpfile, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return Err[*File](err)
	}

	ntmpfile := NewFile(String(tmpfile.Name()))
	ntmpfile.file = tmpfile

	defer ntmpfile.Close()

	return Ok(ntmpfile)
}

// Write opens the named file (creating it with the given permissions if needed),
// then write-locks it and overwrites it with the given content.
func (f *File) Write(content String, mode ...os.FileMode) Result[*File] {
	return f.WriteFromReader(content.Reader(), mode...)
}

// WriteFromReader takes an io.Reader (scr) as input and writes the data from the reader into the file.
// If no FileMode is provided, the default FileMode (0644) is used.
func (f *File) WriteFromReader(scr io.Reader, mode ...os.FileMode) Result[*File] {
	if f.file == nil {
		if mda := f.createAll(); mda.IsErr() {
			return mda
		}
	}

	fmode := FileDefault
	if len(mode) != 0 {
		fmode = mode[0]
	}

	filePath := f.filePath()
	if filePath.IsErr() {
		return Err[*File](filePath.Err())
	}

	f = NewFile(filePath.Ok())

	if err := f.OpenFile(os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fmode).Err(); err != nil {
		return Err[*File](err)
	}

	defer f.Close()

	_, err := io.Copy(f.file, scr)
	if err != nil {
		return Err[*File](err)
	}

	err = f.file.Sync()
	if err != nil {
		return Err[*File](err)
	}

	return Ok(f)
}

// dirPath returns the absolute path of the directory containing the file.
func (f *File) dirPath() Result[String] {
	var (
		path string
		err  error
	)

	if f.IsDir() {
		path, err = filepath.Abs(f.name.Std())
	} else {
		path, err = filepath.Abs(filepath.Dir(f.name.Std()))
	}

	if err != nil {
		return Err[String](err)
	}

	return Ok(String(path))
}

// filePath returns the full file path, including the directory and file name.
func (f *File) filePath() Result[String] {
	dirPath := f.dirPath()
	if dirPath.IsErr() {
		return Err[String](dirPath.Err())
	}

	if f.IsDir() {
		return dirPath
	}

	return Ok(String(filepath.Join(dirPath.Ok().Std(), filepath.Base(f.name.Std()))))
}

func (f *File) createAll() Result[*File] {
	dirPath := f.dirPath()
	if dirPath.IsErr() {
		return Err[*File](dirPath.Err())
	}

	if !f.Exist() {
		if err := os.MkdirAll(dirPath.Ok().Std(), DirDefault); err != nil {
			return Err[*File](err)
		}
	}

	return Ok(f)
}
