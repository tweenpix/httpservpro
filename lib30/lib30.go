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
	ID      int     `db:"id" json:"id"`
	Name    string  `db:"name" json:"name"`
	Age     int     `db:"age" json:"age"`
	Friends []uint8 `db:"friends" json:"friends"`
}

// var users = make(map[int]*User)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var users []User
	query := "SELECT * FROM test_table"
	sql2 := server.Db.Select(&users, query)
	if sql2 != nil {
		panic(sql2.Error())
	}

	for _, user := range users {
		fmt.Fprint(w, "<ul>")
		fmt.Fprintf(w, `<li>ID: %d, Name: %s, Age: %d, Friends: %s</li>`, user.ID, user.Name, user.Age, user.Friends)
		fmt.Fprint(w, "</ul>")
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err2 := server.Db.Exec("INSERT INTO test_table(name, age) VALUES(?, ?)", user.Name, user.Age)
	if err2 != nil {
		fmt.Println(err.Error())
	}

	log.Printf(`User "%v" added to the database!`, user.Name)
}

// 2
func MakeFriends(w http.ResponseWriter, r *http.Request) {
	var friendRequest struct {
		SourceID uint8 `json:"source_id"`
		TargetID uint8 `json:"target_id"`
	}
	var user User
	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	// Получение текущего списка друзей
	// var friends []uint8
	err1 := server.Db.Get(&user, "SELECT friends FROM test_table WHERE id=?", friendRequest.SourceID)
	if err1 != nil {
		log.Fatal(err1)
	}

	// Добавление нового друга в список
	user.Friends = append(user.Friends, uint8(friendRequest.TargetID))

	fmt.Println(user.Friends)

	// Обновление списка друзей в базе данных
	_, err2 := server.Db.Exec("UPDATE test_table SET friends=? WHERE id IN (?,?)", user.Friends, uint8(friendRequest.SourceID), uint8(friendRequest.TargetID))
	if err2 != nil {
		log.Fatal(err2)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s and %s are now friends", friendRequest.TargetID, friendRequest.SourceID)))
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

	// Получение текущего списка друзей
	var friends []int
	err2 := server.Db.Get(&friends, "SELECT friends FROM test_table WHERE id=?", deleteUserRequest.TargetID)
	if err2 != nil {
		return
	}

	// Удаление друга из списка
	for i, id := range friends {
		if id == deleteUserRequest.TargetID {
			friends = append(friends[:i], friends[i+1:]...)
			break
		}
	}

	// Обновление списка друзей в базе данных
	_, err3 := server.Db.Exec("UPDATE test_table SET friends=? WHERE id=?", friends, deleteUserRequest.TargetID)
	if err3 != nil {
		return
	}

	// Удаление пользователя из хранилища
	query := "DELETE FROM test_table WHERE id = ?"

	_, err4 := server.Db.Exec(query, deleteUserRequest.TargetID)
	if err2 != nil {
		fmt.Println(err4.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s has been deleted", deleteUserRequest.TargetID)))
}

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
