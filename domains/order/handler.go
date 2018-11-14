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
}

func NewHandler(s OrderService) Handler {
	return &handler{
		service: s,
	}
}

func (h *handler) HandleUpsert(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var Order Order
	err := decoder.Decode(&Order)
	if err != nil || len(Order.LineItems) <= 0 {
		if err == nil {
			err = errors.New("Order cannot be empty")
		}
		log.Printf("Error handlng upsert: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = h.service.UpsertOrder(&Order); err == nil {
		utils.WriteDataResponse(w, Order)
	} else {
		log.Printf("Failed to create order. Error:\n%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	orderIDstr := mux.Vars(r)["orderid"]
	var order *Order
	if orderID, err := strconv.ParseUint(orderIDstr, 10, 32); err != nil {
		log.Printf("invalid order id: '%d'\n", orderID)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		order = h.service.Get(uint(orderID))
	}
	utils.WriteDataResponse(w, order)
}
