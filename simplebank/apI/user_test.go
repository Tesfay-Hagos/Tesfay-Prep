package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "tesfayprep/simplebank/db/mock"
	db "tesfayprep/simplebank/db/sqlc"
	"tesfayprep/simplebank/util"
	"testing"

	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testcases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			buildStubs: func(store *mockdb.MockStore) {
				//password, _ := util.HashedPassword(user.HashedPassword)
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			//create controller for mockdb
			ctrl := gomock.NewController(t)
			//create instance of store from mockdb
			store := mockdb.NewMockStore(ctrl)
			//tell mockstore to control the function call of mocked methods
			tc.buildStubs(store)
			//create instance of server using given store and test
			server := newTestServer(t, store)
			//create recorder
			recorder := httptest.NewRecorder()
			body := createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			}
			requestbody, _ := json.Marshal(body)
			request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(requestbody))

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckHashedpassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}
func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}
