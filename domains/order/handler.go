package order

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/noahhai/sam-go-ex/utils"
)

type handler struct {
	service OrderService
}

type Handler interface {
	HandleUpsert(w http.ResponseWriter, r *http.Request)
	HandleGet(w http.ResponseWriter, r *http.Request)
	HandleGetAll(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s OrderService) Handler {
	return &handler{
		service: s,
	}
}

func (h *handler) HandleUpsert(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order Order
	err := decoder.Decode(&order)
	if err != nil || (order.ID == 0 && len(order.LineItems) <= 0) {
		if err == nil {
			err = errors.New("Order cannot be empty")
		}
		log.Printf("Error handlng upsert: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = h.service.UpsertOrder(&order); err == nil {
		utils.WriteDataResponse(w, order)
	} else {
		log.Printf("Failed to create order. Error:\n%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	orderIDstr := mux.Vars(r)["orderid"]
	var orders []Order
	if orderID, err := strconv.ParseUint(orderIDstr, 10, 32); err != nil {
		log.Printf("invalid order id: '%d'\n", orderID)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		orders = h.service.Get(uint(orderID))
	}
	utils.WriteDataResponse(w, orders)
}

func (h *handler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	orders := h.service.All()
	utils.WriteDataResponse(w, orders)
}
