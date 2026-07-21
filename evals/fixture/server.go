// HTTP surface for the orders service.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// orderRoutes exposes the order endpoints:
//
//	GET /orders/total?prices=250,100           -> {"total":350}
//	GET /orders/discount?total=350&code=SAVE10 -> {"total":315}
func orderRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders/total", func(w http.ResponseWriter, req *http.Request) {
		var prices []int
		for _, s := range strings.Split(req.URL.Query().Get("prices"), ",") {
			p, err := strconv.Atoi(s)
			if err != nil {
				http.Error(w, "bad price: "+s, http.StatusBadRequest)
				return
			}
			prices = append(prices, p)
		}
		json.NewEncoder(w).Encode(map[string]int{"total": orderTotal(prices)})
	})
	mux.HandleFunc("GET /orders/discount", func(w http.ResponseWriter, req *http.Request) {
		total, err := strconv.Atoi(req.URL.Query().Get("total"))
		if err != nil {
			http.Error(w, "bad total", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]int{"total": applyDiscount(total, req.URL.Query().Get("code"))})
	})
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", orderRoutes()))
}
