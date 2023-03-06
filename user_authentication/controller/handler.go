package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tesfayprep/user_authentication/model"
	"tesfayprep/user_authentication/tokenfunc"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	user := model.UserInfo{}
	json.NewDecoder(r.Body).Decode(&user)
	response := model.Register(user)
	json.NewEncoder(w).Encode(response)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		user := model.UserInfo{}
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Fprintf(w, "invalid body")
			return
		}
		dbuser, isnotuser := model.GetUser(user.Email)
		if !isnotuser {
			if dbuser.Data[0].Password == "" || dbuser.Data[0].Password != user.Password {
				fmt.Fprintf(w, "can not authenticate this user")
				return
			}
			token, err := tokenfunc.GenerateJWT(user.Username)
			if err != nil {
				fmt.Fprintf(w, "error in generating token")
			}

			fmt.Fprintf(w, token)
		}

	case "GET":
		fmt.Fprintf(w, "only POST methods is allowed.")
		return
	}
}
func GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	err := tokenfunc.ValidateToken(w, r)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		users := model.GetAllUser()
		json.NewEncoder(w).Encode(users)
	}
}
