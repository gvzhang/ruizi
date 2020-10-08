package util

import (
	"bytes"
	"encoding/binary"
	"io"
)

func WriteBinary(handle io.Writer, data interface{}) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	_, err = handle.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func ReadBinary(handle io.Reader, readLen uint64, data interface{}) error {
	readData := make([]byte, readLen)
	_, err := handle.Read(readData)
	if err != nil {
		return err
	}

	err = binary.Read(bytes.NewBuffer(readData), binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}
