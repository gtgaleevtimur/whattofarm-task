// Package api - пакет с настройкой API сервера ( роутер,маршруты) + хранилище.
package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

const (
	eu int = 1 + iota
	bu
	both
)

// NewRouter - настраивает роутер и необходимые зависимости.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	data := newData()
	data.refresh()
	control := newController(data)
	r.Post("/api/v1/rates", control.PostPair)
	r.Get("/api/v1/rates", control.GetPair)
	r.NotFound(notFound())
	r.MethodNotAllowed(notAllowed())
	return r
}

// Data - структура хранилища переменных.
type Data struct {
	EthUsdt string `json:"ETH-USDT,omitempty"`
	BtcUsdt string `json:"BTC-USDT,omitempty"`
}

// newData - конструктор хранилища.
func newData() *Data {
	result := &Data{
		EthUsdt: "",
		BtcUsdt: "",
	}
	return result
}

type BinanceNode struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// refresh - обновляет значения переменных согласно курсу бинанс.
func (d *Data) refresh() {
	request, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		log.Fatal(err)
	}
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var temp []BinanceNode
	err = json.Unmarshal(body, &temp)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range temp {
		switch v.Symbol {
		case "ETHUSDT":
			d.EthUsdt = trimZero(v.Price)
		case "BTCUSDT":
			d.BtcUsdt = trimZero(v.Price)
		}
	}
}

// controller - структура контроллера роутера.
type controller struct {
	*Data
}

// newController - конструктор контроллера.
func newController(d *Data) *controller {
	return &controller{
		d,
	}
}

type PostJson struct {
	Pairs []string `json:"pairs"`
}

// PostPair - обработчик POST /api/v1/rates
func (c *controller) PostPair(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var temp PostJson
	err = json.Unmarshal(body, &temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Data.refresh()
	var result Data
	for _, v := range temp.Pairs {
		switch v {
		case "BTC-USDT":
			result.BtcUsdt = c.BtcUsdt
		case "ETH-USDT":
			result.EthUsdt = c.EthUsdt
		}
	}
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetPair - обработчик GET /api/v1/rates
func (c *controller) GetPair(w http.ResponseWriter, r *http.Request) {
	pair := findPair(r.URL.RawQuery)
	c.refresh()
	var temp Data
	switch pair {
	case 0:
		w.WriteHeader(http.StatusBadRequest)
	case 1:
		temp.EthUsdt = c.EthUsdt
		result, err := json.Marshal(temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	case 2:
		temp.BtcUsdt = c.BtcUsdt
		result, err := json.Marshal(temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	case 3:
		temp.BtcUsdt, temp.EthUsdt = c.BtcUsdt, c.EthUsdt
		result, err := json.Marshal(temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}

}

// notFound - обработчик неподдерживаемых маршрутов.
func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

// notAllowed - обработчик неподдерживаемых методов.
func notAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}

// trimZero - обрезает лишние нули значений курсов.
func trimZero(str string) string {
	for strings.HasSuffix(str, "0") {
		temp := []rune(str)
		if temp[len(temp)-3] == '.' {
			return str
		}
		str = strings.TrimSuffix(str, "0")
	}
	return str
}

// findPair - свитчер значений.
func findPair(value string) int {
	temp := strings.Split(value, "=")
	if len(temp) == 1 {
		return 0
	}
	newTemp := strings.Split(temp[1], ",")
	if len(newTemp) == 1 {
		switch newTemp[0] {
		case "ETH-USDT":
			return eu
		case "BTC-USDT":
			return bu
		default:
			return 0
		}
	}
	if len(newTemp) == 2 {
		if (newTemp[0] == "ETH-USDT" && newTemp[1] == "BTC-USDT") || (newTemp[0] == "BTC-USDT" && newTemp[1] == "ETH-USDT") {
			return both
		}
		return 0
	}
	return 0
}
