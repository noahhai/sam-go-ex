package widget

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/noahhai/sam-go-ex/utils"
)

type handler struct {
	service WidgetService
}

type Handler interface {
	HandleUpsert(w http.ResponseWriter, r *http.Request)
	HandleFilter(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s WidgetService) Handler {
	return &handler{
		service: s,
	}
}

func (h *handler) HandleUpsert(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var widget WidgetInventory
	err := decoder.Decode(&widget)
	widget.Category = strings.ToLower(widget.Category)
	widget.Size = strings.ToLower(widget.Size)
	widget.Color = strings.ToLower(widget.Color)
	widget.Name = strings.ToLower(widget.Name)
	if err == nil {
		h.service.UpsertWidget(&widget)
	}
	if err != nil || widget.Name == "" {
		if err == nil {
			err = errors.New("widget name not set")
		}
		log.Printf("Error handlng upsert: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.WriteDataResponse(w, widget)
}

func (h *handler) HandleFilter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	widgets := h.service.Filter(&WidgetFilter{
		Name:     r.FormValue("name"),
		Category: r.FormValue("category"),
		Color:    r.FormValue("color"),
		Size:     r.FormValue("size"),
	})
	utils.WriteDataResponse(w, widgets)
}
