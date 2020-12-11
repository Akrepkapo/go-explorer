package storage

type DbConner interface {
	Initer() error
	Close() error
}

// TransactionStatus is model
type TransactionStatus struct {
	Hash      []byte `gorm:"primary_key;not null"  json:"hash"`
	Params string `gorm:"not null" json:"params"`
	KeyID  string `gorm:"not null" json:"key_id"`
	Time   int64  `gorm:"not null" json:"time"`
	Type   int64  `gorm:"not null" json:"type"`
	Size   int64  `gorm:"not null" json:"size"`

	Ecosystemname string `gorm:"null" json:"ecosystemname"`
	Token_title   string `gorm:"null" json:"token_title"`
	Ecosystem     int64  `gorm:"null" json:"ecosystem"`
}
