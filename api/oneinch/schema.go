package oneinch

// ApproveSpenderResponseDto schema
type ApproveSpenderResponseDto struct {
	Address string `json:"address,omitempty"`
}

// ApproveAllowanceResponseDto schema
type ApproveAllowanceResponseDto struct {
	Allowance string `json:"allowance,omitempty"`
}

// ApproveCalldataResponseDto schema
type ApproveCalldataResponseDto struct {
	Data     string `json:"data,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	To       string `json:"to,omitempty"`
	Value    string `json:"value,omitempty"`
}
