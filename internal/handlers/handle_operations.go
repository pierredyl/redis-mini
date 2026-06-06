package handlers

import (
	"errors"
	"fmt"
	"redis-mini/internal/data"
)

func HandleSet(args []string, data *data.Store) (err error) {
	// Set operation needs at least 3 values, check array size
	if len(args) < 3 {
		return errors.New("error: not enough arguments")
	}

	key := args[1]
	value := args[2]

	data.Set(key, value)

	return nil
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
