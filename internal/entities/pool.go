package entities

type Token struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	WeiVal  int    `json:"wei"`
}

type SwapProtocol struct {
	Name       string `json:"protocolName"`
	Factory    string `json:"factoryAddress"`
	SwapRouter string `json:"swapRouter"`
}

type TokenPair struct {
	Token0 Token `json:"token0"`
	Token1 Token `json:"token1"`
}

type TradePool struct {
	SwapProtocol `json:"tradeProvider"`
	Address      string    `json:"address"`
	Pair         TokenPair `json:"pair"`
}

type TradePair struct {
	Pool0 TradePool `json:"pool0"`
	Pool1 TradePool `json:"pool1"`
}
