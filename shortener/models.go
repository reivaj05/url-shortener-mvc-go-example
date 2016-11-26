package shortener

import (
	"github.com/jinzhu/gorm"
	"github.com/reivaj05/GoJSON"
)

type URLShortModel struct {
	gorm.Model
	LongURL  string `gorm:"type:varchar(100);not null;unique"`
	ShortURL string `gorm:"type:varchar(100);not null;unique"`
}

func (instance *URLShortModel) ToJSON() *GoJSON.JSONWrapper {
	json, _ := GoJSON.New("{}")
	json.SetValueAtPath("id", instance.ID)
	json.SetValueAtPath("long_url", instance.LongURL)
	json.SetValueAtPath("short_url", instance.ShortURL)
	return json
}

var Models = []interface{}{
	&URLShortModel{},
}
