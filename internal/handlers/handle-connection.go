package handlers

import (
	"fmt"
	"net"
	"redis-mini/internal/data"
	"strings"
)

func HandleConnection(conn net.Conn, data *data.Store) (err error) {
	// Defer the close to handle graceful exists
	defer conn.Close()
	args, err := HandleResp(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Checking the first argument for operation
	operation := strings.ToLower(args[0])

	switch operation {
	case "set":
		HandleSet(args, data)
	case "get":
		value := HandleGet(args, data)
		fmt.Println("Got:", value)
	}

	return nil
}
