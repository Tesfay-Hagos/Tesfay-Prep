package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tesfayprep/user_authentication/controller"
	"tesfayprep/user_authentication/model"
	"testing"
	"time"
)

func TestRegiste(t *testing.T) {
	Newuser := model.UserInfo{Username: "Bemnet", Password: "Berut2121", Email: "bemnetthagos@gmail.com", CreatedAt: time.Now()}
	buff := convtobuff(Newuser)
	router := http.NewServeMux()
	router.HandleFunc("/register", controller.RegisterHandler)
	req := httptest.NewRequest(http.MethodPost, "/register", &buff)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := model.JsonResponse{}
	json.NewDecoder(resp.Body).Decode(&response)
	if response.Type != "success" {
		t.Errorf("Test Failed")
	}

}
