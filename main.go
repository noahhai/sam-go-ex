package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/noahhai/sam-go-ex/domains/order"
	"github.com/noahhai/sam-go-ex/domains/widget"

	"log"

	"github.com/apex/gateway"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	flag "github.com/spf13/pflag"
)

const ContentType = "application/json; charset=utf8"

var (
	isServerful = flag.BoolP("serverfull", "s", false, "Should run as lambda, otherwise http server")
)

func RegisterRoutes(db *gorm.DB) {
	widgetService := widget.NewWidgetService(db)
	orderService := order.NewOrderService(db, widgetService)
	widgetHandler := widget.NewHandler(widgetService)
	orderHandler := order.NewHandler(orderService)
	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	r.Use(CorsMiddleware)

	r.HandleFunc("/widgets", widgetHandler.HandleFilter).Methods("GET")
	r.HandleFunc("/widget", widgetHandler.HandleUpsert).Methods("PUT", "POST")
	r.HandleFunc("/order", orderHandler.HandleUpsert).Methods("PUT", "POST")
	r.HandleFunc("/orders", orderHandler.HandleGetAll).Methods("GET")
	r.HandleFunc("/orders/{orderid}", orderHandler.HandleGet).Methods("GET")

	http.Handle("/", r)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	db := getDatabase()
	defer db.Close()
	db.AutoMigrate(&widget.WidgetInventory{}, &order.Order{}, &order.OrderItem{})

	RegisterRoutes(db)
	if !*isServerful {
		log.Println("Starting listening and serving in serverless mode")
		log.Fatal(gateway.ListenAndServe(":5000", nil))
	} else {
		log.Println("Starting listening and serving in serverfull mode")
		log.Fatal(http.ListenAndServe(":5000", nil))
	}
}

func getDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, pass))
	if err != nil {
		panic("failed to connect database")
	}
	return db

}
