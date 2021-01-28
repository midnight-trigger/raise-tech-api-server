package mysql

import (
	"time"
)

type Images struct {
	Id            int64      `gorm:"type:bigint(20);primary_key;auto_increment"`
	ImageFileName string     `gorm:"type:varchar(255);not null"`
	CreatedAt     time.Time  `gorm:"type:datetime;not null"`
	UpdatedAt     time.Time  `gorm:"type:datetime;not null"`
	DeletedAt     *time.Time `gorm:"type:datetime"`
}

//go:generate mockgen -source images.go -destination mock_mysql/mock_images.go
type IImages interface {
	FindById(id int64) (image Images, err error)
	Create(image *Images) (err error)
	Update(oldParams Images, updateParams map[string]interface{}) (err error)
	Delete(image *Images) (err error)
}

func GetNewImage() *Images {
	return &Images{}
}

func (m *Images) FindById(id int64) (image Images, err error) {
	err = db.Where("id = ?", id).First(&image).Error
	return
}

func (m *Images) Create(image *Images) (err error) {
	err = db.Create(&image).Error
	return
}

func (m *Images) Update(oldParams Images, updateParams map[string]interface{}) (err error) {
	err = db.Model(&oldParams).Updates(updateParams).Error
	return
}

func (m *Images) Delete(image *Images) (err error) {
	err = db.Delete(image).Error
	return
}
