package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/service"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RestServer struct {
	cfg         *config.Config
	ctx         context.Context
	Service     service.Service
	UserService service.UserService
}

func NewRestService(ctx context.Context, service service.Service, userService service.UserService) (*RestServer, error) {
	return &RestServer{
		ctx:         ctx,
		cfg:         config.NewConfig(),
		Service:     service,
		UserService: userService,
	}, nil
}

func (s *RestServer) Start() error {
	ctx, stop := signal.NotifyContext(s.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /registration", func(w http.ResponseWriter, r *http.Request) {
		s.Registration(w, r)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		s.Login(w, r)
	})
	mux.HandleFunc("POST /api/item", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.AddItem)
	})
	mux.HandleFunc("GET /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.GetItemById)
	})
	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.GetItems)
	})
	mux.HandleFunc("PUT /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.UpdateItemById)
	})
	mux.HandleFunc("DELETE /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.DeleteItemById)
	})
	mux.HandleFunc("POST /api/shopping_list", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.AddShoppingList)
	})
	mux.HandleFunc("GET /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.GetShoppingListById)
	})
	mux.HandleFunc("GET /api/shopping_lists", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.GetShoppingLists)
	})
	mux.HandleFunc("PUT /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.UpdateShoppingListById)
	})
	mux.HandleFunc("DELETE /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.DeleteShoppingListById)
	})
	mux.HandleFunc("GET /api/shopping_list_items/{id}", func(w http.ResponseWriter, r *http.Request) {
		UserIdentity(w, r, s.GetItemsByShoppingListId)
	})

	serverHTTP := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := serverHTTP.Shutdown(ctxT); err != nil {
		return fmt.Errorf("shutdown server: %s\n", err)
	}
	log.Println("FINAL server shutdown")
	return nil
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
func (s *RestServer) Registration(w http.ResponseWriter, r *http.Request) {
	var user model.CreateUserDTO
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

	uID, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	user.Id = uID.String()
	user.State = 1
	err = s.UserService.CreateUser(r.Context(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true}`)))
	return
}

// Login
// @Summary Логин
// @Tags auths
// @Accept			json
// @Produce		json
// @Param input body model.CreateUserRequest true "Модель которую принимает метод"
// @Success 200 {string}  string "Login successful"
// @Failure 400 {string} string "Invalid request"
// @Router /registration [post]
func (s *RestServer) Login(w http.ResponseWriter, r *http.Request) {
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
	token, err := s.UserService.Login(r.Context(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true, "token": "%s"}`, token)))
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
func (s *RestServer) AddShoppingList(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var sl model.CreateShoppingListDTO
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
	sl.CreatedAt = timestamppb.Now()
	sl.UpdatedAt = timestamppb.Now()
	sl.Items = make([]string, 0)
	sl.State = 2
	err = s.Service.CreateShoppingList(r.Context(), &sl)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
	return
}

// GetShoppingListById
// @Summary Получить список покупок по его id
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Shopping list"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list/{id} [get]
func (s *RestServer) GetShoppingListById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sl, err := s.Service.GetShoppingListById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	data, err := json.Marshal(struct {
		Success      bool `json:"success"`
		ShoppingList *model.ShoppingList
	}{
		Success:      true,
		ShoppingList: sl,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": "Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetShoppingLists
// @Summary Получить все списки покупок
// @Tags shopping_list
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Shopping lists"
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_lists [get]
func (s *RestServer) GetShoppingLists(w http.ResponseWriter, r *http.Request) {
	list, err := s.Service.GetShoppingLists(r.Context())
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	data, err := json.Marshal(struct {
		Success bool `json:"success"`
		Sl      []*model.ShoppingList
	}{
		Success: true,
		Sl:      list,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": "Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
func (s *RestServer) UpdateShoppingListById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var sl *model.UpdateShoppingListDTO
	err = json.Unmarshal(body, &sl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = s.Service.UpdateShoppingList(r.Context(), id, sl)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
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
func (s *RestServer) DeleteShoppingListById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.Service.DeleteShoppingListById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
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
func (s *RestServer) AddItem(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var it model.CreateItemDTO
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
	it.CreatedAt = timestamppb.Now()
	it.UpdatedAt = timestamppb.Now()
	it.IsDone = false
	err = s.Service.CreateItem(r.Context(), &it)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
	return
}

// GetItemById
// @Summary Получить пункт списка покупок по его id
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Item"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/{id} [get]
func (s *RestServer) GetItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := s.Service.GetItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	data, err := json.Marshal(struct {
		Success bool `json:"success"`
		Item    *model.Item
	}{
		Success: true,
		Item:    item,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": "Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetItems
// @Summary Получить все пункты списков покупок
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Items"
// @Failure 400 {string} string "Invalid request"
// @Router /api/items [get]
func (s *RestServer) GetItems(w http.ResponseWriter, r *http.Request) {
	list, err := s.Service.GetItems(r.Context())
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	data, err := json.Marshal(struct {
		Success bool `json:"success"`
		Item    []*model.Item
	}{
		Success: true,
		Item:    list,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": "Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
func (s *RestServer) UpdateItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	var item *model.UpdateItemDTO
	err = json.Unmarshal(body, &item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
		return
	}
	err = s.Service.UpdateItem(r.Context(), id, item)
	if err != nil {
		if errors.Is(err, errors.New("NOT FOUND")) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// DeleteItemById
// @Summary Удалить пункт списка покупок по его id
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}  string "Item deleted"
// @Failure 400 {string} string "Invalid request"
// @Router /api/item/{id} [delete]
func (s *RestServer) DeleteItemById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := s.Service.DeleteItemById(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
}

// GetItemsByShoppingListId
// @Summary Получить список всех покупок по id листа покупок
// @Tags item
// @Accept			json
// @Produce		json
// @Success 200 {string}
// @Failure 400 {string} string "Invalid request"
// @Router /api/shopping_list_items/{id} [get]
func (s *RestServer) GetItemsByShoppingListId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, err := s.Service.GetItemsByShoppingListId(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"success": false, "error": ` + err.Error() + `}`))
			return
		}
	}
	data, err := json.Marshal(struct {
		Success bool `json:"success"`
		Item    []*model.Item
	}{
		Success: true,
		Item:    list,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success": false, "error": "Internal Server Error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
