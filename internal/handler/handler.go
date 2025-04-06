package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

// Login
// @Summary Логин
// @Tags auths
// @Accept			json
// @Produce		json
// @Param input body model.CreateUserRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Login successful"
// @Failure 400 {string} string "Invalid request"
// @Router /registration [post]
func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Username is empty"}`))
		return
	}
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Password is empty"}`))
		return
	}

	token, err := service.Login(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true, "token": "%s"}`, token)))
	return
}

// Registration
// @Summary Регистрация
// @Tags auths
// @Accept			json
// @Produce		json
// @Param input body model.RegistrationUserRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Registration successful"
// @Failure 400 {string} string "Invalid request"
// @Router /login [post]
func Registration(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var user model.User
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Username is empty"}`))
		return
	}
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Password is empty"}`))
		return
	}

	service.Registration(user.Name, user.Password)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true}`)))
	return
}

// AddItem
// @Summary Добавить пункт в список покупок
// @Tags item
// @Accept			json
// @Produce		json
// @Param input body model.CreateItemRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Item added"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/ [post]
func AddItem(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var it model.Item
	err = json.Unmarshal(body, &it)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if it.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Title is empty"}`))
		return
	}
	if it.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "UserId is empty"}`))
		return
	}
	if it.ShoppingListId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "ShoppingListId is empty"}`))
		return
	}
	slID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	it.Id = slID.String()
	it.CreatedAt = time.Now()
	it.UpdatedAt = time.Now()
	it.IsDone = false
	repository.CheckInterface(&it)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
	return
}

// AddShoppingList
// @Summary Создать список покупок
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Param input body model.CreateShoppingListRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Shopping list added"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list/ [post]
func AddShoppingList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var sl model.ShoppingList
	err = json.Unmarshal(body, &sl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if sl.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Title is empty"}`))
		return
	}
	if sl.UserId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "UserId is empty"}`))
		return
	}
	slID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	sl.Id = slID.String()
	sl.CreatedAt = time.Now()
	sl.UpdatedAt = time.Now()
	sl.Items = make([]string, 0)
	sl.State = 1
	repository.CheckInterface(&sl)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
	return
}

// GetItems
// @Summary Получить все пункты списков покупок
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Items"
// @Failure 400 {string} string "Invalid request"
// @Router /api/items [get]
func GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	list := repository.GetItems()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "items": ` + list + `}`))
}

// GetShoppingLists
// @Summary Получить все списки покупок
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Shopping lists"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_lists [get]
func GetShoppingLists(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	list := repository.GetSls()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "shopping_list": ` + list + `}`))
}

// GetItemById
// @Summary Получить пункт списка покупок по его id
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Item"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/{id} [get]
func GetItemById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := repository.GetItemById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Item with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "item": ` + item + `}`))
}

// GetShoppingListById
// @Summary Получить список покупок по его id
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Shopping list"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list/{id} [get]
func GetShoppingListById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sl, err := repository.GetSlById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Shopping list with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "shopping_list": ` + sl + `}`))
}

// DeleteItemById
// @Summary Удалить пункт списка покупок по его id
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Item deleted"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/{id} [delete]
func DeleteItemById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := repository.DeleteItemById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Item with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// DeleteShoppingListById
// @Summary Удалить список покупок по его id
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Shopping list deleted"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list/{id} [delete]
func DeleteShoppingListById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := repository.DeleteSlById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Shopping list with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// UpdateShoppingListById
// @Summary Обновить список покупок по его id
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Param input body model.UpdateShoppingListRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Shopping list updated"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list/{id} [put]
func UpdateShoppingListById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = repository.UpdateSl(id, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Shopping list with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// UpdateItemById
// @Summary Обновить пункт списка покупок по его id
// @Tags item
// @Accept			json
// @Produce		json
// @Param input body model.UpdateItemRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Item updated"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/{id} [put]
func UpdateItemById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = repository.UpdateItem(id, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": "Item with this Id not found"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}
