package order

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/noahhai/sam-go-ex/domains/widget"
)

type OrderService interface {
	Get(orderID uint) []Order
	All() []Order
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

func (s *orderService) All() []Order {
	var orders []Order
	s.db.Preload("LineItems").Preload("LineItems.Widget").Find(&orders)
	return orders
}

func (s *orderService) Get(orderID uint) []Order {
	var orders []Order
	s.db.Preload("LineItems").Preload("LineItems.Widget").First(&orders, orderID)
	return orders
}

func (s *orderService) UpsertOrder(o *Order) error {
	var currOrder *Order
	if o.ID != 0 {
		var currOrders []Order
		s.db.Preload("LineItems").Preload("LineItems.Widget").First(&currOrders, o.ID)
		if currOrders != nil && len(currOrders) > 0 {
			currOrder = &currOrders[0]
		}
	}

	numWidgets := len(o.LineItems)

	// space accounts for products on current order that have been removed
	inventoryChanges := make(map[uint]int, 3*numWidgets)
	productsInOrder := make([]*widget.WidgetInventory, 0, numWidgets)

	for i, li := range o.LineItems {
		inventory := s.ws.Get(li.WidgetID)
		if inventory == nil {
			return fmt.Errorf("could not find widget for line item; id:'%s', attr:'%v'", li.WidgetID, li.Widget)
		}
		productsInOrder = append(productsInOrder, inventory)
		inventoryChanges[li.WidgetID] = -li.Quantity

		if li.Quantity == 0 {
			o.LineItems = append(o.LineItems[:i], o.LineItems[i+1:]...)
		}
	}
	if currOrder != nil {
		for _, li := range currOrder.LineItems {
			_, ok := inventoryChanges[li.WidgetID]
			if !ok {
				productsInOrder = append(productsInOrder, &li.Widget)
			}
			inventoryChanges[li.WidgetID] += li.Quantity
		}
	}

	for _, inventory := range productsInOrder {
		if -inventoryChanges[inventory.ID] > inventory.Remaining {
			return fmt.Errorf("only %d left in inventory for widget; id:'%d'", inventory.Remaining, inventory.ID)
		}
		inventory.Remaining = inventory.Remaining + inventoryChanges[inventory.ID]
	}
	for _, inventory := range productsInOrder {
		s.ws.UpsertWidget(inventory)
	}

	if currOrder == nil {
		s.db.Create(o)
	} else {
		s.db.Delete(OrderItem{}, "order_id = ?", o.ID)
		s.db.Save(o)
	}
	return nil
}
