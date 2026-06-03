package handlers

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) (err error) {
	// Defer the close to handle graceful exists
	defer conn.Close()
	args, err := HandleResp(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Arguments:", args)

	return nil
}
