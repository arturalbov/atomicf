package atomicf

import (
	"crypto/sha256"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const logFilePostfix = "-log-*.alog"

// Atomic write on linux https://danluu.com/file-consistency/
type AtomicFile struct {
	*os.File
}

func OpenFile(name string, flag int, perm os.FileMode) (*AtomicFile, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &AtomicFile{file}, nil
}

// todo: for now working only with one log file at time
func (file *AtomicFile) Recover() (err error) {
	logFiles, err := filepath.Glob(filepath.Join(file.dir(), file.logFilePattern()))
	if err != nil || len(logFiles) == 0 {
		return
	}

	// take only first file
	data, err := ioutil.ReadFile(logFiles[0])
	if err != nil {
		return
	}

	// means that log file not corrupted. otherwise do nothing and delete file
	if len(data) > 32 && VerifyHash(data[0:32], data[32:]) {
		offset := binary.LittleEndian.Uint64(data[32:40])
		_, err = file.writePostLog(logFiles[0], data[40:], int64(offset))
		return
	}

	if err = os.Remove(logFiles[0]); err != nil {
		return
	}

	return
}

func (file *AtomicFile) Write(p []byte) (n int, err error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return
	}

	logFileName, err := file.logOperation(p, fileInfo.Size())
	if err != nil {
		return
	}

	var dir *os.File
	if dir, err = os.Open(file.dir()); err != nil {
		return
	}
	defer dir.Close()

	if err = dir.Sync(); err != nil {
		return
	}

	if n, err = file.File.Write(p); err != nil {
		return
	}

	if err = file.File.Sync(); err != nil {
		return
	}

	if err = os.Remove(logFileName); err != nil {
		return
	}

	if err = dir.Sync(); err != nil {
		return
	}

	return
}

func (file *AtomicFile) WriteAt(b []byte, off int64) (n int, err error) {
	return file.writeLog(b, off)
}

func (file *AtomicFile) writeLog(b []byte, off int64) (n int, err error) {
	logFileName, err := file.logOperation(b, off)
	if err != nil {
		return
	}

	return file.writePostLog(logFileName, b, off)
}

func (file *AtomicFile) writePostLog(logFileName string, data []byte, offset int64) (n int, err error) {
	var dir *os.File
	if dir, err = os.Open(file.dir()); err != nil {
		return
	}
	defer dir.Close()

	if err = dir.Sync(); err != nil {
		return
	}

	if n, err = file.File.WriteAt(data, offset); err != nil {
		return
	}

	if err = file.File.Sync(); err != nil {
		return
	}

	if err = os.Remove(logFileName); err != nil {
		return
	}

	if err = dir.Sync(); err != nil {
		return
	}

	return
}

// return directory of current file
func (file *AtomicFile) dir() string {
	return filepath.Dir(file.File.Name())
}

func (file *AtomicFile) logFilePattern() string {
	fileName := filepath.Base(file.Name())

	// remove dot postfix
	dotIndex := strings.Index(fileName, ".")
	if dotIndex != -1 {
		fileName = fileName[0:dotIndex]
	}

	return fileName + logFilePostfix
}

// return path of create log file
func (file *AtomicFile) logOperation(b []byte, off int64) (path string, err error) {
	offsetBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(offsetBytes, uint64(off))
	data := append(offsetBytes, b...)
	dataHash := sha256.Sum256(data)

	return WriteTempFile(
		file.dir(),
		file.logFilePattern(),
		append(dataHash[:], data...),
	)
}
