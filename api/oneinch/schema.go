package oneinch

type TokenAddress string

type TokenBalance string

// ApproveSpenderResponseDto schema
type ApproveSpenderResponseDto struct {
	// Address of the 1inch router that must be trusted to spend funds for the exchange
	Address string `json:"address,omitempty"`
}

// ApproveAllowanceParamsDto schema
type ApproveAllowanceParamsDto struct {
	TokenAddress  string `json:"tokenAddress,omitempty" url:"tokenAddress"`
	WalletAddress string `json:"walletAddress,omitempty" url:"walletAddress"`
}

// ApproveAllowanceResponseDto schema
type ApproveAllowanceResponseDto struct {
	Allowance string `json:"allowance,omitempty"`
}

// ApproveCalldataParamsDto schema
type ApproveCalldataParamsDto struct {
	TokenAddress string `json:"tokenAddress,omitempty" url:"tokenAddress"`
	Amount       string `json:"amount,omitempty" url:"amount"`
}

// ApproveCalldataResponseDto schema
type ApproveCalldataResponseDto struct {
	// The encoded data to call the approve method on the swapped token contract
	Data string `json:"data,omitempty"`
	// Gas price for fast transaction processing
	GasPrice string `json:"gasPrice,omitempty"`
	// Token address that will be allowed to exchange through 1inch router
	To string `json:"to,omitempty"`
	// Native token value in WEI (for approve is always 0)
	Value string `json:"value,omitempty"`
}

// ProtocolsResponseDto schema
type ProtocolsResponseDto struct {
	// List of protocols that are available for routing in the 1inch Aggregation protocol
	Protocols []ProtocolImageDto `json:"protocols,omitempty"`
}

// ProtocolImageDto schema
type ProtocolImageDto struct {
	// Protocol id
	ID string `json:"id,omitempty"`
	// Protocol title
	Title string `json:"title,omitempty"`
	// Protocol logo image
	Img string `json:"img,omitempty"`
}

// TokensResponseDto schema
type TokensResponseDto struct {
	// List of supported tokens
	Tokens map[string]TokenDto `json:"tokens,omitempty"`
}

// TokenDto schema
type TokenDto struct {
	// Symbol for the token
	Symbol string `json:"symbol,omitempty"`
	// Name of the token
	Name string `json:"name,omitempty"`
	// Address of the token
	Address string `json:"address,omitempty"`
	// Number of decimal places for the token
	Decimals uint64 `json:"decimals"`
	// URL for the logo of the token
	LogoURI string `json:"logoURI,omitempty"`
}

// PresetsResponseDto schema
type PresetsResponseDto struct {
	MaxResult []PresetDto `json:"MAX_RESULT,omitempty"`
	LowestGas []PresetDto `json:"LOWEST_GAS,omitempty"`
}

// PresetDto schema
type PresetDto struct {
	ComplexityLevel uint64 `json:"complexityLevel,omitempty"`
	MainRouterParts uint64 `json:"mainRouteParts,omitempty"`
	Parts           uint64 `json:"parts,omitempty"`
	VirtualParts    uint64 `json:"virtualParts,omitempty"`
}

// TransactionDto schema
type TransactionDto struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Data     string `json:"data,omitempty"`
	Value    string `json:"value,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Gas      string `json:"gas,omitempty"`
}

// PathViewDto schema
type PathViewDto struct {
	Name             string `json:"name,omitempty"`
	Part             string `json:"part,omitempty"`
	FromTokenAddress string `json:"fromTokenAddress,omitempty"`
	ToTokenAddress   string `json:"toTokenAddress,omitempty"`
}

// QuoteParamsDto schema
type QuoteParamsDto struct {
	FromTokenAddress string `json:"fromTokenAddress,omitempty" url:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress,omitempty" url:"toTokenAddress"`
	Amount           string `json:"amount,omitempty" url:"amount"`
}

// QuoteResponseDto schema
type QuoteResponseDto struct {
	FromToken       TokenDto `json:"fromToken,omitempty"`
	ToToken         TokenDto `json:"toToken,omitempty"`
	FromTokenAmount string   `json:"fromTokenAmount,omitempty"`
	ToTokenAmount   string   `json:"toTokenAmount,omitempty"`
	EstimatedGas    uint64   `json:"estimatedGas,omitempty"`
}

// SwapParamsDto schema
type SwapParamsDto struct {
	FromTokenAddress string  `json:"fromTokenAddress,omitempty" url:"fromTokenAddress"`
	ToTokenAddress   string  `json:"toTokenAddress,omitempty" url:"toTokenAddress"`
	Amount           string  `json:"amount,omitempty" url:"amount"`
	FromAddress      string  `json:"fromAddress,omitempty" url:"fromAddress"`
	Slippage         float64 `json:"slippage,omitempty" url:"slippage"`
	DestReceiver     string  `json:"destReceiver,omitempty" url:"destReceiver"`
	Fee              float64 `json:"fee,omitempty" url:"fee"`
	GasPrice         string  `json:"gasPrice,omitempty" url:"gasPrice"`
	ComplexityLevel  uint64  `json:"complexityLevel,string,omitempty" url:"complexityLevel"`
	ConnectorTokens  uint64  `json:"connectorTokens,string,omitempty" url:"connectorTokens"`
	AllowPartialFill bool    `json:"allowPartialFill,omitempty" url:"allowPartialFill"`
	DisableEstimate  bool    `json:"disableEstimate,omitempty" url:"disableEstimate"`
	GasLimit         string  `json:"gasLimit,omitempty" url:"gasLimit"`
	MainRouteParts   uint64  `json:"mainRouteParts,string,omitempty" url:"mainRouteParts"`
	Parts            uint64  `json:"parts,string,omitempty" url:"parts"`
}

// SwapResponseDto schema
type SwapResponseDto struct {
	FromToken       TokenDto `json:"fromToken,omitempty"`
	ToToken         TokenDto `json:"toToken,omitempty"`
	FromTokenAmount string   `json:"fromTokenAmount,omitempty"`
	ToTokenAmount   string   `json:"toTokenAmount,omitempty"`
	From            string   `json:"from,omitempty"`
	To              string   `json:"to,omitempty"`
	Data            string   `json:"data,omitempty"`
	Value           string   `json:"value,omitempty"`
	GasPrice        string   `json:"gasPrice,omitempty"`
	Gas             string   `json:"gas,omitempty"`
}
