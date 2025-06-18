package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/model"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var TOKEN string

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("HTTP API Консольный клиент")
	fmt.Println("---------------------")
	fmt.Println("Доступные команды: get-lists, get-items, create-list, create-item, delete-list, delete-item, exit")
	if TOKEN == "" {
		fmt.Println("\nУ вас уже есть аакаунт?")
		fmt.Print("Выберете действие: ")
		fmt.Println("1. Регистрация")
		fmt.Println("2. Логин")
		fmt.Print("-> ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Введите имя: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Введите пароль: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)
			registration(name, password)
		case "2":
			fmt.Print("Введите имя: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Введите пароль: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)
			login(name, password)
		default:
			fmt.Println("Неверный выбор")
		}

	}
	for {
		fmt.Print("\n-> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "get-lists":
			getSls()
		case "get-items":
			getItems()
		case "create-list":
			fmt.Print("Введите название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			fmt.Print("Введите userId: ")
			userId, _ := reader.ReadString('\n')
			userId = strings.TrimSpace(userId)

			createSl(title, userId)
		case "create-item":
			fmt.Print("Введите название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)
			fmt.Print("Введите userId: ")
			userId, _ := reader.ReadString('\n')
			userId = strings.TrimSpace(userId)
			fmt.Print("Введите shoppingListId: ")
			shoppingListId, _ := reader.ReadString('\n')
			shoppingListId = strings.TrimSpace(shoppingListId)

			createItem(title, userId, shoppingListId)

		case "delete-list":
			fmt.Print("Введите id: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			deleteList(id)
		case "delete-item":
			fmt.Print("Введите id: ")
			id, _ := reader.ReadString('\n')
			id = strings.TrimSpace(id)

			deleteItem(id)
		case "exit":
			fmt.Println("Выход из программы...")
			return
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func login(name, password string) {
	createUserRequest := model.CreateUserRequest{
		Name:     name,
		Password: password,
	}
	jsonData, err := json.Marshal(createUserRequest)
	if err != nil {
		fmt.Println("Ошибка при создании JSON:", err)
		return
	}
	resp, err := http.Post(
		"http://0.0.0.0:8989/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	var U *model.LoginUserResponse
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &U)
	TOKEN = U.Token
}

func registration(name, password string) {
	createUserRequest := model.CreateUserRequest{
		Name:     name,
		Password: password,
	}
	jsonData, err := json.Marshal(createUserRequest)
	if err != nil {
		fmt.Println("Ошибка при создании JSON:", err)
		return
	}
	_, err = http.Post(
		"http://0.0.0.0:8989/registration",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	login(name, password)
}

func getSls() {
	req, err := http.NewRequest("GET", "http://0.0.0.0:8989/api/shopping_lists", nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}
	req.Header.Add("Authorization", "Bearer "+TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	type S struct {
		Success string `json:"success,omitempty"`
		SL      []*model.ShoppingList
	}
	var s S
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)
	for _, l := range s.SL {
		fmt.Printf(l.String())
	}
}

func getItems() {
	req, err := http.NewRequest("GET", "http://0.0.0.0:8989/api/items", nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}
	req.Header.Add("Authorization", "Bearer "+TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	type S struct {
		Success string `json:"success,omitempty"`
		Item    []*model.Item
	}
	var s S
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &s)
	for _, i := range s.Item {
		fmt.Printf(i.String())
	}
}

func createSl(title, userId string) {
	type S struct {
		Title  string `json:"title,omitempty"`
		UserId string `json:"user_id,omitempty"`
	}
	createSlRequest := S{
		Title:  title,
		UserId: userId,
	}
	jsonData, err := json.Marshal(createSlRequest)
	if err != nil {
		fmt.Println("Ошибка при создании JSON:", err)
		return
	}
	req, err := http.NewRequest("POST",
		"http://0.0.0.0:8989/api/shopping_list",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	var r Response
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	fmt.Println(r.String())
}

func createItem(title, userId, shoppingListId string) {
	type S struct {
		Title          string `json:"title,omitempty"`
		UserId         string `json:"user_id,omitempty"`
		ShoppingListId string `json:"shopping_list_id,omitempty"`
	}
	createSlRequest := S{
		Title:          title,
		UserId:         userId,
		ShoppingListId: shoppingListId,
	}
	jsonData, err := json.Marshal(createSlRequest)
	if err != nil {
		fmt.Println("Ошибка при создании JSON:", err)
		return
	}
	req, err := http.NewRequest("POST",
		"http://0.0.0.0:8989/api/item",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	var r Response
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	fmt.Println(r.String())
}

func deleteList(id string) {
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("http://0.0.0.0:8989/api/shopping_list/%s", id),
		nil,
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	var r Response
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	fmt.Println(r.String())
}

func deleteItem(id string) {
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("http://0.0.0.0:8989/api/item/%s", id),
		nil,
	)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()
	var r Response
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &r)
	fmt.Println(r.String())
}

type Response struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r Response) String() string {
	return fmt.Sprintf("Success: \"%t\", Error: \"%s\"", r.Success, r.Error)
}
