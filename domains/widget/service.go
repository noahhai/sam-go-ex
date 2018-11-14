package widget

import (
	"github.com/jinzhu/gorm"
)

type WidgetService interface {
	Filter(f *WidgetFilter) []WidgetInventory
	Get(widgetID uint) *WidgetInventory
	UpsertWidget(w *WidgetInventory)
	CreateWidget(w *WidgetInventory)
}

type widgetService struct {
	db *gorm.DB
}

func NewWidgetService(db *gorm.DB) WidgetService {
	return &widgetService{
		db: db,
	}
}

func (s *widgetService) Get(widgetID uint) (widget *WidgetInventory) {
	var widgets []WidgetInventory
	s.db.First(&widgets, widgetID)
	if widgets != nil && len(widgets) > 0 {
		widget = &(widgets[0])
	}
	return widget
}

func (s *widgetService) Filter(f *WidgetFilter) []WidgetInventory {
	widgets := []WidgetInventory{}
	s.db.Where(f).Find(&widgets)
	return widgets
}

func (s *widgetService) CreateWidget(w *WidgetInventory) {
	s.db.Create(w)
}

func (s *widgetService) UpsertWidget(w *WidgetInventory) {
	// new widget but check if already a record for widget
	// with these attrs
	if w.ID == 0 {
		filter := WidgetFilter{
			Name:     w.Name,
			Category: w.Category,
			Color:    w.Color,
			Size:     w.Size,
		}
		var currWidgets []WidgetInventory

		// first seems to require reference to slice
		s.db.Where(&filter).First(&currWidgets)
		if currWidgets != nil && len(currWidgets) >= 1 {
			currWidget := currWidgets[0]
			currWidget.Remaining = w.Remaining
			*w = currWidget
		}
	}
	s.db.Save(w)
}
