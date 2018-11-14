package widget

import (
	"github.com/jinzhu/gorm"
)

type WidgetFilter struct {
	Name     string
	Category string
	Color    string
	Size     string
}

func (WidgetFilter) TableName() string {
	return "widget_inventories"
}

type WidgetInventory struct {
	gorm.Model
	Name      string
	Category  string
	Color     string
	Size      string
	Remaining int
}
