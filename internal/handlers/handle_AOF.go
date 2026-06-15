package handlers

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strconv"
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

func HandleAOFRead() (err error) {
	file, err := os.OpenFile("database.aof", os.O_RDONLY, 0644)
	if err != nil {
		return errors.New("failed to open AOF file")
	}
	defer file.Close()

	// Start reading the file line by line
	reader := bufio.NewReader(file)

	nextByte, err := reader.ReadByte()
	if err != nil {
		return errors.New("error: reading from AOF")
	}

	switch nextByte {
	case '*':
		// The start of a new array
		err := ParseAOFArray(reader)
		if err != nil {
			return errors.New("error: parsing array from AOF")
		}
	}

	return nil
}

func ParseAOFArray(reader *bufio.Reader) (err error) {
	// Get the array size
	arraySizeStr, err := readLine(reader)
	if err != nil {
		return errors.New("error: failed getting array size from AOF")
	}

	arraySize, err := strconv.Atoi(arraySizeStr)
	if err != nil {
		return errors.New("error: failed converting AOF array size string to int")
	}

	for i := 0; i < arraySize; i++ {
		nextByte, err := reader.ReadByte()
		if err != nil {
			return errors.New("error: reading next byte from AOF array")
		}

		switch nextByte {
		case '$':
			// Process a bulk string
		}
	}

	return nil
}
