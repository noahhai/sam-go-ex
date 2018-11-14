package order

import (
	"github.com/jinzhu/gorm"
	"github.com/noahhai/sam-go-ex/domains/widget"
)

type Order struct {
	gorm.Model
	LineItems []OrderItem `gorm:"foreignkey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	WidgetID uint
	OrderID  uint
	Widget   widget.WidgetInventory `gorm:"foreignkey:WidgetID"`
	Quantity int
}
