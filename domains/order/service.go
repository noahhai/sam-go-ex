package order

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/noahhai/sam-go-ex/domains/widget"
)

type OrderService interface {
	Get(orderID uint) *Order
	UpsertOrder(w *Order) error
}

type orderService struct {
	db *gorm.DB
	ws widget.WidgetService
}

func NewOrderService(db *gorm.DB, ws widget.WidgetService) OrderService {
	return &orderService{
		db: db,
		ws: ws,
	}
}

func (s *orderService) Get(orderID uint) (order *Order) {
	var orders []Order
	s.db.Preload("LineItems").Preload("LineItems.Widget").First(&orders, orderID)
	if orders != nil && len(orders) > 0 {
		order = &(orders[0])
	}
	return order
}

func (s *orderService) UpsertOrder(o *Order) error {
	var currOrder *Order
	if o.ID != 0 {
		var currOrders []Order
		s.db.First(&currOrders, o.ID)
		if currOrders != nil && len(currOrders) > 0 {
			currOrder = &currOrders[0]
		}
	}

	numWidgets := len(o.LineItems)
	inventoryChanges := make(map[uint]int, numWidgets)
	productsInOrder := make([]*widget.WidgetInventory, 0, numWidgets)

	for _, li := range o.LineItems {
		inventory := s.ws.Get(li.WidgetID)
		if inventory == nil {
			return fmt.Errorf("could not find widget for line item; id:'%s', attr:'%v'", li.WidgetID, li.Widget)
		}
		productsInOrder = append(productsInOrder, inventory)
		inventoryChanges[li.WidgetID] = li.Quantity
		if currOrder != nil {
			var currLineItem OrderItem
			for _, c := range currOrder.LineItems {
				if c.WidgetID == li.WidgetID {
					currLineItem = c
				}
			}
			if currLineItem.ID != 0 {
				inventoryChanges[li.WidgetID] -= currLineItem.Quantity
			}
		}
		for _, inventory := range productsInOrder {
			if inventoryChanges[inventory.ID] > inventory.Remaining {
				return fmt.Errorf("only %d left in inventory for widget; id:'%s'", inventory.Remaining, inventory.ID)
			}
			inventory.Remaining -= inventoryChanges[inventory.ID]
		}
		for _, inventory := range productsInOrder {
			s.ws.UpsertWidget(inventory)
		}
	}

	if currOrder == nil {
		s.db.Create(o)
	} else {
		s.db.Save(o)
	}
	return nil
}
