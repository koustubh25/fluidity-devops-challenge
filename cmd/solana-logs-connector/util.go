package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func websocketWriteSubscribeMessage(w io.WriteCloser) (err error) {
	defer w.Close()

	err = json.NewEncoder(w).Encode(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "logsSubscribe",
		"params": []map[string]any{
			{
				"mentions": []string{"11111111111111111111111111111111"},
			},
		},
	})

	if err != nil {
		return fmt.Errorf(
			"Failed to write to the websocket the hello message! %v",
			err,
		)
	}

	return
}
