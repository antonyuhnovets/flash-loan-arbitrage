package entities

type Token struct {
	Name    string `json:"name" bson:"name" swaggertype:"string"`
	Address string `json:"address" bson:"address" swaggertype:"string"`
	WeiVal  int    `json:"wei" bson:"wei" swaggertype:"integer"`
}

type SwapProtocol struct {
	Name       string `json:"protocolName" bson:"protocolName" swaggertype:"string"`
	Factory    string `json:"factoryAddress" bson:"factoryAddress" swaggertype:"string"`
	SwapRouter string `json:"swapRouter" bson:"swapRouter" swaggertype:"string"`
}

type TokenPair struct {
	Token0 Token `json:"token0" bson:"token0"`
	Token1 Token `json:"token1" bson:"token1"`
}

type TradePool struct {
	SwapProtocol `json:"tradeProvider"`
	Address      string    `json:"address" bson:"address" swaggertype:"string"`
	Pair         TokenPair `json:"pair" bson:"pair"`
}

type TradePair struct {
	Pool0 TradePool `json:"pool0" bson:"pool0"`
	Pool1 TradePool `json:"pool1" bson:"pool1"`
}
