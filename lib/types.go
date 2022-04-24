package devops_hiring_challenge

type (
	// SolanaWebsocketLog is received off the connection from Solana
	SolanaWebsocketLog struct {
		Result *int `json:"result"`
		Params struct {
			Result struct {
				Value struct {
					Logs []string `json:"logs"`
				} `json:"value"`
			} `json:"result"`
		} `json:"params"`
	}

	// ComputeUnits is sent down Kafka with the number
	ComputeUnits uint64
)
