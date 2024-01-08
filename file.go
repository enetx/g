package g

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"gitlab.com/x0xO/g/internal/filelock"
)

// NewFile returns a new File instance with the given name.
func NewFile(name String) *File { return &File{name: name} }

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

// Close closes the File and unlocks its underlying file, if it is not already closed.
func (f *File) Close() error {
	if f.file == nil {
		return fmt.Errorf("%s: file is already closed and unlocked", f.name)
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

// ReadLines reads the file and returns its content as a slice of lines.
func (f *File) ReadLines() Result[Slice[String]] {
	read := f.Read()
	if read.IsErr() {
		return Err[Slice[String]](read.Err())
	}

	return Ok(read.Ok().SplitLines())
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
		return Err[*File](fmt.Errorf("no such file: %s", f.name))
	}

	nf := NewFile(newpath).createAll()

	if err := os.Rename(f.name.Std(), newpath.Std()); err != nil {
		return Err[*File](err)
	}

	return nf
}

// SeekToLine moves the file pointer to the specified line number and reads the
// specified number of lines from that position.
// The function returns the new position and a concatenation of the lines read as an String.
//
// Parameters:
//
// - position int64: The starting position in the file to read from
//
// - linesRead int: The number of lines to read.
//
// Returns:
//
// - int64: The new file position after reading the specified number of lines
//
// - String: A concatenation of the lines read as an String.
//
// Example usage:
//
//	f := g.NewFile("file.txt")
//	position, content := f.SeekToLine(0, 5) // Read 5 lines from the beginning of the file
func (f *File) SeekToLine(position int64, linesRead int) (int64, String) {
	iterator := f.Iter().Unwrap()

	if _, err := f.file.Seek(position, 0); err != nil {
		f.Close()
		return 0, ""
	}

	var (
		content     strings.Builder
		linesReaded int
	)

	for line := iterator.ByLines(); line.Next(); linesReaded++ {
		if linesReaded == linesRead {
			f.Close()
			break
		}

		_, _ = content.WriteString(line.ToString().Add("\n").Std())
		position += int64(line.ToBytes().Len() + 1) // Add 1 for the newline character
	}

	return position, String(content.String())
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
		return ToResult(f.file.Stat())
	}

	return ToResult(os.Stat(f.name.Std()))
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

	stat := f.Stat()

	if stat.IsOk() && stat.Ok().IsDir() {
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

	stat := f.Stat()
	if stat.IsOk() && stat.Ok().IsDir() {
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

// Iter returns a new fiter instance that can be used to read the file
// line by line, word by word, rune by rune, or byte by byte.
//
// Returns:
//
// - *fiter: A pointer to the new fiter instance.
//
// Example usage:
//
//	f := g.NewFile("file.txt")
//	iterator := f.Iter().Unwrap() // Returns a new fiter instance for the file
func (f *File) Iter() Result[*fiter] {
	openfile := f.Open()
	if openfile.IsErr() {
		return Err[*fiter](openfile.Err())
	}

	openfile.Ok().fiter = &fiter{
		scanner: bufio.NewScanner(f.file),
		file:    f,
	}

	return Ok(f.fiter)
}

// Buffer sets the initial buffer to use when iterating and the maximum size of the buffer.
//
// By default, Iterator uses an internal buffer and grows it as large as necessary.
// This method allows you to use a custom buffer and limit its size.
//
// Parameters:
//
// - buf: A byte slice that will be used as a buffer.
//
// - max: The maximum size of the buffer.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByRunes()
//
//	customBuffer := make([]byte, 1024)
//	iterator.Buffer(customBuffer, 4096)
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.ToString())
//	}
func (fit fiter) Buffer(buf []byte, max int) { fit.scanner.Buffer(buf, max) }

// By configures the fiter instance's scanner to use a custom split function.
//
// The custom split function should take a byte slice and a boolean indicating whether this is the
// end of file.
// It should return the advance count, the token, and any encountered error.
// Don't forget to close the file in the custom split function!
//
// Parameters:
//
// - f: A split function of the form func(data []byte, atEOF bool) (advance int, token []byte, err
// error).
//
// Returns:
//
// - An fiter instance with the scanner configured to use the provided custom split function.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	customSplitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
//	    // Custom implementation here
//	}
//
//	iterator := myFile.Iter().Unwrap().By(customSplitFunc)
//
//	for iterator.Next() {
//	    fmt.Printf("%s", iterator.ToString())
//	}
func (fit fiter) By(f func(data []byte, atEOF bool) (advance int, token []byte, err error)) fiter {
	fit.scanner.Split(bufio.SplitFunc(f))
	return fit
}

// Bytes sets the iterator to read the file byte by byte.
//
// Returns:
//
// - An fiter instance with the scanner configured to read the file byte by byte.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByBytes()
//
//	for iterator.Next() {
//	    fmt.Println(iterator.ToString())
//	}
func (fit fiter) ByBytes() fiter {
	fit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			fit.file.Close()
			return 0, nil, nil
		}

		return 1, data[0:1], nil
	})

	return fit
}

// Lines sets the iterator to read the file line by line.
//
// Returns:
//
// - An fiter instance with the scanner configured to read the file line by line.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByLines()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.ToBytes())
//	}
func (fit fiter) ByLines() fiter {
	dropCR := func(data []byte) []byte {
		if len(data) > 0 && data[len(data)-1] == '\r' {
			return data[0 : len(data)-1]
		}

		return data
	}

	fit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			fit.file.Close()
			return 0, nil, nil
		}

		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, dropCR(data[0:i]), nil
		}

		if atEOF {
			return len(data), dropCR(data), nil
		}

		return 0, nil, nil
	})

	return fit
}

// Runes sets the iterator to read the file rune by rune.
//
// Returns:
//
// - An fiter instance with the scanner configured to read the file rune by rune.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByRunes()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.ToString())
//	}
func (fit fiter) ByRunes() fiter {
	fit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			fit.file.Close()
			return 0, nil, nil
		}

		if data[0] < utf8.RuneSelf {
			return 1, data[0:1], nil
		}

		_, width := utf8.DecodeRune(data)
		if width > 1 {
			return width, data[0:width], nil
		}

		if !atEOF && !utf8.FullRune(data) {
			return 0, nil, nil
		}

		return 1, []byte(string(utf8.RuneError)), nil
	})

	return fit
}

// Words sets the iterator to read the file word by word.
//
// Returns:
//
// - An fiter instance with the scanner configured to read the file word by word.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByWords()
//
//	for iterator.Next() {
//	    fmt.Println(iterator.ToString())
//	}
func (fit fiter) ByWords() fiter {
	fit.By(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			fit.file.Close()
			return 0, nil, nil
		}

		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if !unicode.IsSpace(r) {
				break
			}
		}

		for width, i := 0, start; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if unicode.IsSpace(r) {
				return i + width, data[start:i], nil
			}
		}

		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}

		return start, nil, nil
	})

	return fit
}

// Error returns the first non-EOF error encountered by the Iterator.
//
// Call this method after an iteration loop has finished to check if any errors occurred
// during the iteration process.
//
// Returns:
//
// - An error encountered by the iterator, or nil if no errors occurred.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByRunes()
//
//	for iterator.Next() {
//	    fmt.Printf("%c", iterator.ToString())
//	}
//
//	if err := iterator.Err(); err != nil {
//	    log.Printf("Error while iterating: %v", err)
//	}
func (fit fiter) Err() error { return fit.scanner.Err() }

// Next advances the iterator to the next item (byte, line, word, or rune) and
// returns true if successful or false if there are no more items to read.
//
// Returns:
//
// - A boolean value indicating whether the iterator successfully advanced to the next item.
// Returns false if there are no more items to read.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap()
//	lines := iterator.ByLines() // or iterator.ByWords() or iterator.ByRunes()
//
//	for lines.Next() {
//	    fmt.Println(iterator.ToString())
//	}
func (fit fiter) Next() bool { return fit.scanner.Scan() }

// ToBytes returns the current item as an Bytes instance.
//
// Returns:
//
// - An Bytes instance containing the current item in the iterator.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByBytes() // Sets the iterator to read the file byte by byte
//
//	for iterator.Next() {
//	    fmt.Println(iterator.ToBytes())
//	}
func (fit fiter) ToBytes() Bytes { return Bytes(fit.scanner.Bytes()) }

// String returns the current item as an String instance.
//
// Returns:
//
// - An String instance containing the current item in the iterator.
//
// Example usage:
//
//	myFile := g.NewFile("path/to/myfile.txt")
//
//	iterator := myFile.Iter().Unwrap().ByLines() // Sets the iterator to read the file line by line
//
//	for iterator.Next() {
//	    fmt.Println(iterator.ToString())
//	}
func (fit fiter) ToString() String { return String(fit.scanner.Text()) }
