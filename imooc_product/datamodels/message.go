package datamodels

// 消息结构体
type Message struct {
	ProductID int64
	UserID    int64
}

// 创建结构体
func NewMessage(userID int64, productID int64) *Message {
	return &Message{ProductID:productID, UserID:userID}
}