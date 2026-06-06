package handlers

import (
	"bytes"
	"errors"
	"os"
)

func HandleAOFWrite(buffer *bytes.Buffer) (err error) {
	// Open file
	file, err := os.OpenFile("database.aof", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("failed to open AOF file")
	}
	defer file.Close()

	// Append the buffer
	if _, err := file.Write(buffer.Bytes()); err != nil {
		return errors.New("failed to write buffer to AOF")
	}

	return nil
}
