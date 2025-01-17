package model

// RealEstate 房产信息
type RealEstate struct {
	ID string `gorm:"primaryKey" json:"id"` // 房产ID
}

// Transaction 交易信息
type Transaction struct {
	ID           string `gorm:"primaryKey" json:"id"`      // 交易ID
	RealEstateID string `json:"realEstateId" gorm:"index"` // 房产ID
}

// TableName 指定表名
func (RealEstate) TableName() string {
	return "real_estates"
}

// TableName 指定表名
func (Transaction) TableName() string {
	return "transactions"
}
