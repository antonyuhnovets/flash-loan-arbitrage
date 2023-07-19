package entities

type Token struct {
	ID      int    `json:"id" bson:"id" gorm:"column:id;primaryKey;type:integer;autoIncrement:true"`
	Name    string `json:"name" bson:"name" gorm:"column:name;type:varchar(40)"`
	Address string `json:"address" bson:"address" gorm:"column:address;type:varchar(50)"`
	Wei     int    `json:"wei" bson:"wei" gorm:"column:wei;type:bigint"`
}

type SwapProtocol struct {
	ID         int    `json:"id" bson:"id" gorm:"column:id;primaryKey;type:integer;autoIncrement:true"`
	Name       string `json:"name" bson:"name" gorm:"column:name;type:varchar(40)"`
	Factory    string `json:"factory" bson:"factory" gorm:"column:factory;type:varchar(50)"`
	SwapRouter string `json:"router" bson:"router" gorm:"column:router;type:varchar(50)"`
}

type TokenPair struct {
	ID       int   `json:"id" bson:"id" gorm:"column:id;primaryKey;type:integer;autoIncrement:true"`
	Token0   Token `json:"token0" bson:"token0" gorm:"foreignKey:Token0ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Token0ID int   `json:"-"`
	Token1   Token `json:"token1" bson:"token1" gorm:"foreignKey:Token1ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Token1ID int   `json:"-"`
}

type Pool struct {
	ID         int          `json:"id" bson:"id" gorm:"column:id;primaryKey;type:integer;autoIncrement:true;"`
	Address    string       `json:"address" bson:"address" gorm:"column:address;type:varchar(50)"`
	Pair       TokenPair    `json:"pair" bson:"pair" gorm:"foreignKey:PairID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PairID     int          `json:"-"`
	Protocol   SwapProtocol `json:"protocol" bson:"protocol" gorm:"foreignKey:ProtocolID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ProtocolID int          `json:"-"`
}

type TradePair struct {
	Pool0 Pool `json:"pool0"`
	Pool1 Pool `json:"pool1"`
}
