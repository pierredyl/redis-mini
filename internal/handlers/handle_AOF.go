package handlers

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"redis-mini/internal/data"
	"strings"
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

func HandleAOFRead(data *data.Store) (err error) {
	file, err := os.OpenFile("database.aof", os.O_RDONLY, 0644)
	if err != nil {
		return errors.New("failed to open AOF file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		args, _, err := HandleResp(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		operation := strings.ToLower(args[0])
		switch operation {
		case "set":
			err := HandleSet(args, data)
			if err != nil {
				return errors.New("Set failed")
			}
		}
	}
	return nil
}
