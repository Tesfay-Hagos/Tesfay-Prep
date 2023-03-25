package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	mockdb "tesfayprep/simplebank/db/mock"
	db "tesfayprep/simplebank/db/sqlc"
	"tesfayprep/simplebank/util"
	"tesfayprep/token"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenmaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currence: util.RandomCurrency(),
	}
}
func CreateAccount(t *testing.T) {
	user := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenmaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "createaccountOK",
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)

			},
			buildStubs: func(store *mockdb.MockStore) {
				createaccountparams := db.CreateaccountParams{Owner: account.Owner, Balance: 0, Currence: account.Currence}
				/*store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil) */
				store.EXPECT().Createaccount(gomock.Any(),
					gomock.Eq(createaccountparams)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//request
			accountReq := createAccountRequest{Currency: account.Currence}

			// serialize the request body to JSON
			body, err := json.Marshal(accountReq)
			if err != nil {
				// handle error
				log.Fatalf("error marshaling requestparam:%s", err)
			}

			// create a new HTTP request with the serialized JSON body
			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(body))
			if err != nil {
				log.Fatalf("error creating request:%s", err)
			}
			request.Header.Set("Content-Type", "application/json")

			//request
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

/*
func TestListAccount(t *testing.T) {
	user := randomUser(t)
	account := randomAccount(user.Username)
	accounts := []db.Account{}
	accounts = append(accounts, account)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenmaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "listaccountOK",
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				ListAccountsParams := db.ListAccountsParams{Owner: account.Owner, Limit: 5, Offset: 0}
				/*store.EXPECT().
				GetAccount(gomock.Any(), gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil) */
/*
				store.EXPECT().ListAccounts(gomock.Any(),
					gomock.Eq(ListAccountsParams)).Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			header := http.Header{}
			header.Add("page_id", "1")
			header.Add("page_size", "5")

			// create a new HTTP request with the serialized JSON body
			request, err := http.NewRequest(http.MethodGet, "/accounts", nil)
			if err != nil {
				log.Fatalf("error creating request:%s", err)
			}
			request.Header = header

			//request
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
*/
// check the account in the body and initial account created for test
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func randomUser(t *testing.T) (user db.User) {
	user = db.User{
		Username: util.RandomOwner(),
		FullName: util.RandomOwner(),
	}
	return
}
