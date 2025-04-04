package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	_ "github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

func Base(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "data": "Main page"}`))
	return
}

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

	userID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	user.State = 1

	secretKey := []byte(config.Cfg.Secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID.String(),
		"name":  user.Name,
		"state": 1,
	})
	tokenString, err := token.SignedString(secretKey)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true, "token": "%s"}`, tokenString)))
	return
}

func Add(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if strings.Contains(string(body), "shopping_list_id") {
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
		it.Id = slID.String()
		it.CreatedAt = time.Now()
		it.UpdatedAt = time.Now()
		it.IsDone = false
		repository.CheckInterface(&it)

	} else {
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
		sl.Id = slID.String()
		sl.CreatedAt = time.Now()
		sl.UpdatedAt = time.Now()
		sl.Items = make([]string, 0)
		sl.State = 1
		repository.CheckInterface(&sl)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
	return
}

func GetItems(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	list := repository.GetItems()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "items": ` + list + `}`))
}

func GetSls(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	list := repository.GetSls()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "shopping_list": ` + list + `}`))
}

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

func GetSlById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func DeleteSlById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func UpdateSlById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
