package atomicf

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicFile_CreateWrite(t *testing.T) {
	testData := []byte{0, 1, 2, 3, 4, 5, 6, 7}

	dir, err := ioutil.TempDir("", "atomic-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fileName := filepath.Join(dir, "atomic-test")
	atomic, err := OpenFile(fileName, os.O_CREATE|os.O_SYNC|os.O_WRONLY, 0666)

	if err != nil {
		t.Errorf("error open atomic file %s: %v", fileName, err)
		return
	}

	if _, err = atomic.Write(testData); err != nil {
		t.Errorf("error writing to atomic file: %v", err)
	}

	if err = atomic.Close(); err != nil {
		t.Errorf("error closing atomic file: %v", err)
	}

	writtenData, err := ioutil.ReadFile(atomic.Name())

	if err != nil {
		t.Errorf("error reading atomic file: %v", err)
	}

	if bytes.Compare(writtenData, testData) != 0 {
		t.Errorf("writen bytes incorrect. Expected: %v; Got: %v", testData, writtenData)
	}
}

func TestAtomicFile_WriteAt(t *testing.T) {

	fileData := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	changeData := []byte{8, 9}
	changeOffset := int64(2)
	resultData := []byte{0, 1, 8, 9, 4, 5, 6, 7}

	dir, err := ioutil.TempDir("", "atomic-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fileName, err := WriteTempFile(dir, "atomic-test-*", fileData)
	if err != nil {
		t.Fatalf("error creating test atomic file: %v", err)
	}

	atomic, err := OpenFile(fileName, os.O_CREATE|os.O_SYNC|os.O_WRONLY, 0666)
	if err != nil {
		t.Errorf("error open atomic file %s: %v", fileName, err)
		return
	}

	_, err = atomic.WriteAt(changeData, changeOffset)
	if err != nil {
		t.Fatalf("error writing to atomic file: %v", err)
	}

	if err = atomic.Close(); err != nil {
		t.Errorf("error closing atomic file: %v", err)
	}

	writtenData, err := ioutil.ReadFile(atomic.Name())

	if err != nil {
		t.Errorf("error reading atomic file: %v", err)
	}

	if bytes.Compare(writtenData, resultData) != 0 {
		t.Errorf("writen bytes incorrect. Expected: %v; Got: %v", resultData, writtenData)
	}
}

func TestAtomicFile_Recover(t *testing.T) {

	offsetBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(offsetBytes, uint64(2))
	data := append(offsetBytes, 8, 9)
	dataHash := sha256.Sum256(data)
	logFileData := append(dataHash[:], data...)

	fileData := []byte{0, 1, 2, 3, 4, 5, 6, 7}

	resultData := []byte{0, 1, 8, 9, 4, 5, 6, 7}

	atomicFileName := "atomic-test"

	dir, err := ioutil.TempDir("", "atomic-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	logFileName, err := WriteTempFile(dir, atomicFileName+logFilePostfix, logFileData)
	if err != nil {
		t.Fatalf("error creating test log file: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(dir, atomicFileName), fileData, 0666)
	if err != nil {
		t.Fatalf("error creating test atomic file: %v", err)
	}

	atomic, err := OpenFile(filepath.Join(dir, atomicFileName), os.O_CREATE|os.O_SYNC|os.O_WRONLY, 0666)
	if err != nil {
		t.Errorf("error open atomic file %s: %v", logFileName, err)
		return
	}

	err = atomic.Recover()
	if err != nil {
		t.Fatalf("error recovering atomic file: %v", err)
	}

	if err = atomic.Close(); err != nil {
		t.Errorf("error closing atomic file: %v", err)
	}

	writtenData, err := ioutil.ReadFile(atomic.Name())

	if err != nil {
		t.Errorf("error reading atomic file: %v", err)
	}

	if bytes.Compare(writtenData, resultData) != 0 {
		t.Errorf("writen bytes incorrect. Expected: %v; Got: %v", resultData, writtenData)
	}
}
