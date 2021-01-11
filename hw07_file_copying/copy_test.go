package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const outFileName = "out.txt"

func TestCopy(t *testing.T) {
	defer os.Remove(outFileName)

	t.Run("src file not exists err", func(t *testing.T) {
		err := Copy("random_file_name.txt", outFileName, 0, 0)
		require.Error(t, err)
	})

	t.Run("src file is not regular", func(t *testing.T) {
		err := Copy("/dev/urandom", outFileName, 0, 0)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("offset exceeds src file size", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile("", "copy_test")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString("qwer")
		if err != nil {
			log.Fatal(err)
		}

		err = Copy(tmpFile.Name(), outFileName, 5, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})

	srcFileData := []byte{'q', 'w', 'e'}
	successTestCases := []struct {
		limit               int64
		offset              int64
		srcFileData         []byte
		expectedOutFileData []byte
	}{
		{0, 0, srcFileData, []byte{'q', 'w', 'e'}},
		{1000, 0, srcFileData, []byte{'q', 'w', 'e'}},
		{1, 0, srcFileData, []byte{'q'}},
		{2, 0, srcFileData, []byte{'q', 'w'}},
		{0, 1, srcFileData, []byte{'w', 'e'}},
		{0, 1, srcFileData, []byte{'w', 'e'}},
		{0, 2, srcFileData, []byte{'e'}},
		{0, 3, srcFileData, []byte{}},
		{1, 1, srcFileData, []byte{'w'}},
		{2, 2, srcFileData, []byte{'e'}},
	}

	for _, tt := range successTestCases {
		t.Run(fmt.Sprintf("success Copy with limit=%d,offset=%d", tt.limit, tt.offset), func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "copy_test")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.Write(tt.srcFileData)
			if err != nil {
				log.Fatal(err)
			}

			err = Copy(tmpFile.Name(), outFileName, tt.offset, tt.limit)
			require.NoError(t, err)

			outFile, err := os.Open(outFileName)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err = outFile.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			outFileData, err := ioutil.ReadFile(outFileName)
			if err != nil {
				log.Fatal(err)
			}

			require.Equal(t, tt.expectedOutFileData, outFileData)
		})
	}
}

func TestAdjustLimit(t *testing.T) {

	tests := []struct {
		name        string
		limitArg    int64
		fileSize    int64
		offset      int64
		expectedRes int64
	}{
		{"limit not defined or zero case", 0, 666, 0, 666},
		{"limit exceeds bytes to write size ", 667, 666, 0, 666},
		{"limit exceeds bytes to write size, with offset", 567, 666, 100, 566},
		{"ok case", 3, 666, 0, 3},
		{"ok case, with offset", 3, 666, 333, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := adjustLimit(tt.limitArg, tt.fileSize, tt.offset)
			require.Equal(t, tt.expectedRes, res)
		})
	}
}
