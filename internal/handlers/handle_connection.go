package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"redis-mini/internal/data"
	"strings"
)

func HandleConnection(conn net.Conn, data *data.Store) (err error) {
	// Defer the close to handle graceful exists
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		args, buffer, err := HandleResp(reader)
		fmt.Println(args)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Checking the first argument for operation
		operation := strings.ToLower(args[0])

		switch operation {
		case "set":
			err := HandleSet(args, data)
			if err != nil {
				conn.Write([]byte("-ERR set failed\r\n"))
				return errors.New("Set failed")
			}

			// Use Buffer to write to AOF on successful set
			err = HandleAOFWrite(buffer)
			if err != nil {
				conn.Write([]byte("-ERR AOF write failed\r\n"))
				return errors.New("Write to AOF failed")
			}

			// Set was successful, write back
			conn.Write([]byte("+OK\r\n"))

		case "get":
			value, ok := HandleGet(args, data)
			// Get was not successful
			if !ok {
				conn.Write([]byte("$-1\r\n"))
			} else {
				// Get was successful. Write back to the CLI
				value_string := value.(string)
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value_string), value_string)))
			}

		case "ping":
			conn.Write([]byte("+PONG\r\n"))

		case "command":
			// Currently, not supporting any special commands
			conn.Write([]byte("*0\r\n"))

		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
