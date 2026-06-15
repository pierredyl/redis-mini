package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func HandleResp(conn net.Conn) ([]string, *bytes.Buffer, error) {
	var buffer bytes.Buffer
	var args []string
	reader := bufio.NewReader(conn)

	// The very first byte in RESP indicates type
	typeByte, err := reader.ReadByte()
	if err != nil {
		return nil, nil, err
	}

	// Append typeByte to the buffer
	buffer.WriteByte(typeByte)

	// If the type is an array
	if typeByte == '*' {
		args, err = ParseArray(reader, &buffer, args)
		if err != nil {
			fmt.Println("error: failed ParseArray", err)
			return nil, nil, err
		}
	}

	// Print the buffer
	//fmt.Println("buffer contents:", buffer)
	// Print as a direct string
	//fmt.Println("buffer as string:", buffer.String())

	return args, &buffer, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	// Reads until the end of CRLF
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\r\n"), nil

}

func ParseArray(reader *bufio.Reader, buffer *bytes.Buffer, args []string) ([]string, error) {
	lengthStr, err := readLine(reader)
	if err != nil {
		fmt.Println("Failed to read array length", err)
		return nil, err
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		fmt.Println("Could not convert string to integer", err)
		return nil, err
	}

	buffer.Write([]byte(lengthStr))
	buffer.Write([]byte("\r\n"))

	for i := 0; i < int(length); i++ {
		nextByte, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}

		// Bulk string
		if nextByte == '$' {
			buffer.WriteByte(nextByte)
			args, err = ParseBulkString(reader, buffer, args)
			if err != nil {
				fmt.Println("Failed to parse bulk string", err)
				return nil, err
			}

		}
	}
	return args, nil
}

func ParseBulkString(reader *bufio.Reader, buffer *bytes.Buffer, args []string) ([]string, error) {
	bulkStringLengthStr, err := readLine(reader)
	if err != nil {
		fmt.Println("Failed to get bulk string length", err)
		return nil, err
	}

	bulkStringLength, err := strconv.Atoi(bulkStringLengthStr)
	if err != nil {
		fmt.Println("error: converting bulkStringLengthStr to int", err)
		return nil, err
	}

	buffer.Write([]byte(bulkStringLengthStr))
	buffer.Write([]byte("\r\n"))

	var bulkString []byte = make([]byte, bulkStringLength)

	_, err = io.ReadFull(reader, bulkString)
	if err != nil {
		fmt.Println("Error reading bulk string", err)
		return nil, err
	}

	buffer.Write(bulkString)
	buffer.Write([]byte("\r\n"))

	// Discard the escape bytes \r\n
	_, err = reader.Discard(2)
	if err != nil {
		fmt.Println("Failed to discard escape bytes parsing bulk string", err)
		return nil, err
	}

	args = append(args, string(bulkString))

	return args, nil
}
