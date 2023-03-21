package lib30

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"package30/server"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []uint8 `json:"friends"`
}

// var users = make(map[int]*User)

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
	// user.ID = nextID
	// nextID++

	// users[user.ID] = &user

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
	// fmt.Print(jsonData)
	// log.Printf("%v - %d", user.Name, user.Age)

	_, err2 := server.Db.Exec("INSERT INTO test_table(name, age) VALUES(?, ?)", user.Name, user.Age)
	if err2 != nil {
		fmt.Println(err.Error())
	}

	log.Printf(`User "%v" added to the database!`, user.Name)
}

// // 2
// func MakeFriends(w http.ResponseWriter, r *http.Request) {
// 	var friendRequest struct {
// 		SourceID int `json:"source_id"`
// 		TargetID int `json:"target_id"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&friendRequest)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	sourceUser, sourceExists := users[friendRequest.SourceID]
// 	targetUser, targetExists := users[friendRequest.TargetID]

// 	if !sourceExists || !targetExists {
// 		http.Error(w, "one or both users do not exist", http.StatusBadRequest)
// 		return
// 	}

// 	sourceUser.Friends = append(sourceUser.Friends, targetUser.Name)
// 	targetUser.Friends = append(targetUser.Friends, sourceUser.Name)

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(fmt.Sprintf("%s and %s are now friends", sourceUser.Name, targetUser.Name)))
// }

// // 3
// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	var deleteUserRequest struct {
// 		TargetID int `json:"target_id"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&deleteUserRequest)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	user, exists := users[deleteUserRequest.TargetID]
// 	if !exists {
// 		http.Error(w, "user does not exist", http.StatusBadRequest)
// 		return
// 	}

// 	// Удаление пользователя из хранилища
// 	delete(users, user.ID)

// 	// Удаление пользователя из списков друзей всех его друзей
// 	for _, friendName := range user.Friends {
// 		friend, friendExists := findUserByName(friendName)
// 		if friendExists {
// 			friend.Friends = remove(friend.Friends, user.Name)
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(fmt.Sprintf("%s has been deleted", user.Name)))
// }

// // функция поиска пользователя по имени
// func findUserByName(name string) (*User, bool) {
// 	for _, user := range users {
// 		if user.Name == name {
// 			return user, true
// 		}
// 	}
// 	return nil, false
// }

// // функция удаления элемента из слайса
// func remove(slice []string, element string) []string {
// 	for i, s := range slice {
// 		if s == element {
// 			return append(slice[:i], slice[i+1:]...)
// 		}
// 	}
// 	return slice
// }

// // 4
// func GetUserFriends(w http.ResponseWriter, r *http.Request) {
// 	// Извлекаем ID пользователя из URL-адреса
// 	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	//fmt.Fprint(w, userID)
// 	user, exists := users[userID]
// 	if !exists {
// 		http.Error(w, "user does not exist", http.StatusBadRequest)
// 		return
// 	}

// 	response := make(map[string]interface{})
// 	response["user"] = user.Name
// 	response["friends"] = user.Friends

// 	// Отправляем ответ в формате JSON
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(response)
// }

// 5
func UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	var users []User

	// Извлекаем ID пользователя из URL-адреса
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT * FROM test_table WHERE id = ?"
	sql2 := server.Db.Select(&users, query, userID)
	if sql2 != nil {
		panic(sql2.Error())
	}

	for _, user := range users {

		fmt.Printf("ID: %d, Name: %s, Age: %d, Friends: %s\n", user.ID, user.Name, user.Age, user.Friends)
	}

	// Декодируем JSON-тело запроса
	var requestBody map[string]string
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Извлекаем новый возраст пользователя из JSON-тела запроса
	newAge, err := strconv.Atoi(requestBody["age"])
	fmt.Printf("Устанавливаем возраст: %d\n", newAge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем возраст пользователя в мапе

	query_update := "UPDATE test_table SET age = ? WHERE id = ?"

	_, err2 := server.Db.Exec(query_update, newAge, userID)
	if err2 != nil {
		fmt.Println(err.Error())
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "возраст пользователя успешно обновлён")

}
