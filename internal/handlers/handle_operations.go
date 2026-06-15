package handlers

import (
	"errors"
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

func HandleGet(args []string, data *data.Store) (interface{}, bool) {
	// Get Operations need at least 2 values
	if len(args) < 2 {
		return nil, false
	}

	key := args[1]

	value, ok := data.Get(key)

	return value, ok

}
