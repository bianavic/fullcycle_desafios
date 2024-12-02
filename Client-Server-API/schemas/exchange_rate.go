package schemas

import "gorm.io/gorm"

type Rate struct {
	gorm.Model
	ID  int64  `gorm:"primaryKey;autoIncrement"`
	Bid string `json:"bid"`
}
