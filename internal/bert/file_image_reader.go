package bert

import (
	"bytes"
	"io"
	"os"
)

type fileImageReader struct {
	path string
	file *os.File
}

func newFileImageReader(inputPath string) ImageReader {
	return &fileImageReader{inputPath, nil}
}

func newBytesImageReader(source []byte) ImageReader {
	return io.NopCloser(bytes.NewReader(source))
}

func (rd *fileImageReader) Read(b []byte) (int, error) {
	if rd.file == nil {
		f, err := os.Open(rd.path)

		if err != nil {
			return -1, err
		}
		rd.file = f
	}

	return rd.file.Read(b)
}

func (rd *fileImageReader) Close() error {
	if rd.file != nil {
		return rd.file.Close()
	}
	return nil
}
