package atomicf

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
)

// returns temp file name
func WriteTempFile(dir string, pattern string, data []byte) (path string, err error) {
	tmp, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return
	}

	n, err := tmp.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}

	if err1 := tmp.Sync(); err == nil {
		err = err1
	}

	if err1 := tmp.Close(); err == nil {
		err = err1
	}

	return tmp.Name(), nil
}

func VerifyHash(hashBytes []byte, data []byte) bool {
	dataHashBytes := sha256.Sum256(data)
	return bytes.Compare(hashBytes, dataHashBytes[:]) == 0
}
