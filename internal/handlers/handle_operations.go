package handlers

import (
	"fmt"
	"redis-mini/internal/data"
)

func HandleSet(args []string, data *data.Store) {
	// Set operation needs at least 3 values, check array size
	if len(args) < 3 {
		fmt.Println("error: not enough arguments in set")
		return
	}

	key := args[1]
	value := args[2]

	data.Set(key, value)
}

func HandleGet(args []string, data *data.Store) interface{} {
	// Get Operations need at least 2 values
	if len(args) < 2 {
		fmt.Println("error: not enough arguments in get")
		return nil
	}

	key := args[1]

	value, ok := data.Get(key)
	if !ok {
		fmt.Println("error: failed to retrieve value with key:", key)
		return nil
	}

	return value

}
