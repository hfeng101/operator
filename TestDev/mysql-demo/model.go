package mysql_demo

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TestTable struct {
	Name     string `json:"name:omitempty,type:varchar(128)"`
	Password string `json:"password:omitempty,type:varchar(128)"`
}

type ResourceTable struct {
}

type User struct {
	ID       uint   `json:"primaryKey"`
	Name     string `gorm:"size:100"`
	Age      int
	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt gorm.DeleteAt `gorm:"index"`
}
