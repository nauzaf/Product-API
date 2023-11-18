package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	app *App
}

func newHttpHandler() *httpHandler {
	app := newApp()
	return &httpHandler{
		app: app,
	}
}

func (h *httpHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.app.getProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if products == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(products)
}

func (h *httpHandler) getProduct(w http.ResponseWriter, r *http.Request) {
	idStr, err := url.PathUnescape(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	intVar, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.app.getProduct(intVar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(product)
}

func (h *httpHandler) postProduct(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, ok := requestBody["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isExist, err := h.app.isProductNameExist(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isExist {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = h.app.createProduct(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := struct {
		Message string `json:"message"`
	}{
		Message: "product created successfully",
	}
	json.NewEncoder(w).Encode(respBody)
}

func (h *httpHandler) increaseInventory(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productIDStr, ok := requestBody["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := h.app.getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.app.updateInventory(productID, p.Inventory+1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := struct {
		Message string `json:"message"`
	}{
		Message: "product inventory increased successfully",
	}
	json.NewEncoder(w).Encode(respBody)
}

func (h *httpHandler) decreaseInventory(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productIDStr, ok := requestBody["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := h.app.getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p.Inventory == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if p.Expiry != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.app.updateInventory(productID, p.Inventory-1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := struct {
		Message string `json:"message"`
	}{
		Message: "product inventory decreased successfully",
	}
	json.NewEncoder(w).Encode(respBody)
}

func (h *httpHandler) setExpired(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productIDStr, ok := requestBody["product_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.app.setProductExpired(productID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := struct {
		Message string `json:"message"`
	}{
		Message: "product successfully set as expired",
	}
	json.NewEncoder(w).Encode(respBody)
}
