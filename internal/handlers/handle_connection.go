package handlers

import (
	"errors"
	"fmt"
	"net"
	"redis-mini/internal/data"
	"strings"
)

func HandleConnection(conn net.Conn, data *data.Store) (err error) {
	// Defer the close to handle graceful exists
	defer conn.Close()
	args, buffer, err := HandleResp(conn)
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
			return errors.New("Set failed")
		}
		// Use Buffer to write to AOF on successful set
		err = HandleAOFWrite(buffer)
		if err != nil {
			return errors.New("Write to AOF failed")
		}

	case "get":
		value := HandleGet(args, data)
		fmt.Println("Got:", value)
	}

	return nil
}
