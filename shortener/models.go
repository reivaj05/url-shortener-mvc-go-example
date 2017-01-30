package shortener

import (
	"time"

	"github.com/reivaj05/GoJSON"
)

type URLShortModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	LongURL   string `gorm:"type:varchar(100);not null;unique"`
	ShortURL  string `gorm:"type:varchar(100);not null;unique"`
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
