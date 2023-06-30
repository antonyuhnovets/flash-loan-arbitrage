package entities

type Token struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
	WeiVal  int    `json:"wei" bson:"wei"`
}

type SwapProtocol struct {
	Name       string `json:"protocolName"`
	Factory    string `json:"factoryAddress"`
	SwapRouter string `json:"swapRouter"`
}

type TokenPair struct {
	Token0 Token `json:"token0" bson:"token0"`
	Token1 Token `json:"token1" bson:"token1"`
}

type TradePool struct {
	Address  string       `json:"address"`
	Pair     TokenPair    `json:"pair"`
	Protocol SwapProtocol `json:"protocol"`
}

type TradePair struct {
	Pool0 TradePool `json:"pool0" bson:"pool0"`
	Pool1 TradePool `json:"pool1" bson:"pool1"`
}
