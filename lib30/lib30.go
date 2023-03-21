package lib30

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

var users = make(map[int]*User)
var nextID = 1
var Database *sql.DB

func main() {
	// подключение к базе данных
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/test")
	Database = db
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `
hello world:)
	`)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Fprint(w, &user)
	user.ID = nextID
	nextID++

	users[user.ID] = &user

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonData := json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
	fmt.Print(jsonData)

	log.Printf("%v - %d", user.Name, user.Age)

	insert, err := Database.Query("INSERT INTO test_table(name, age) VALUES(?, ?)", user.Name, user.Age)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer insert.Close()

	log.Println("User added to the database!")
}

// 2
func MakeFriends(w http.ResponseWriter, r *http.Request) {
	var friendRequest struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sourceUser, sourceExists := users[friendRequest.SourceID]
	targetUser, targetExists := users[friendRequest.TargetID]

	if !sourceExists || !targetExists {
		http.Error(w, "one or both users do not exist", http.StatusBadRequest)
		return
	}

	sourceUser.Friends = append(sourceUser.Friends, targetUser.Name)
	targetUser.Friends = append(targetUser.Friends, sourceUser.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s and %s are now friends", sourceUser.Name, targetUser.Name)))
}

// 3
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var deleteUserRequest struct {
		TargetID int `json:"target_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&deleteUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[deleteUserRequest.TargetID]
	if !exists {
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из хранилища
	delete(users, user.ID)

	// Удаление пользователя из списков друзей всех его друзей
	for _, friendName := range user.Friends {
		friend, friendExists := findUserByName(friendName)
		if friendExists {
			friend.Friends = remove(friend.Friends, user.Name)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s has been deleted", user.Name)))
}

// функция поиска пользователя по имени
func findUserByName(name string) (*User, bool) {
	for _, user := range users {
		if user.Name == name {
			return user, true
		}
	}
	return nil, false
}

// функция удаления элемента из слайса
func remove(slice []string, element string) []string {
	for i, s := range slice {
		if s == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// 4
func GetUserFriends(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID пользователя из URL-адреса
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Fprint(w, userID)
	user, exists := users[userID]
	if !exists {
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	response := make(map[string]interface{})
	response["user"] = user.Name
	response["friends"] = user.Friends

	// Отправляем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 5
func UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID пользователя из URL-адреса
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, exists := users[userID]
	if !exists {
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}

	// Декодируем JSON-тело запроса
	var requestBody map[string]string
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Извлекаем новый возраст пользователя из JSON-тела запроса
	newAge, err := strconv.Atoi(requestBody["new age"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем возраст пользователя в мапе
	user.Age = newAge
	users[userID] = user

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "возраст пользователя успешно обновлён")
}
