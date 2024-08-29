package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/sherifzaher/clone-simplebank/db/mock"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  0,
		Currency: util.RandomCurrency(),
	}
}

func TestAccountAPI(t *testing.T) {
	account := randomAccount(util.RandomOwner())

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore)
		buildRequest  func() (string, string, io.Reader)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Create Account",
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
					}).
					Times(1).
					Return(account, nil)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := fmt.Sprintf("/accounts")
				arg := CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
				}
				params, err := json.Marshal(arg)
				require.NoError(t, err)
				return url, http.MethodPost, bytes.NewReader(params)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(store)

			server := newTestStore(t, store)

			url, method, body := tc.buildRequest()
			request, err := http.NewRequest(method, url, body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//
	//store := mockdb.NewMockStore(ctrl)
	//server := newTestServer(t, store)
	//
	//storeArg := db.CreateAccountParams{
	//	Owner:    createdAccount.Owner,
	//	Currency: createdAccount.Currency,
	//	Balance:  0,
	//}
	//
	//store.
	//	EXPECT().
	//	CreateAccount(gomock.Any(), storeArg).
	//	Times(1).
	//	Return(db.Account{}, nil)
	//
	//params := CreateAccountParams{
	//	Owner:    createdAccount.Owner,
	//	Currency: createdAccount.Currency,
	//}
	//
	//arg, err := json.Marshal(params)
	//require.NoError(t, err)
	//
	//url := fmt.Sprintf("/accounts")
	//request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(arg))
	//recorder := httptest.NewRecorder()
	//require.NoError(t, err)
	//
	//server.router.ServeHTTP(recorder, request)
	//checkResponse(t, recorder)
}

func checkResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	fmt.Println(recorder.Body, recorder.Code)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Currency, gotAccount.Currency)
}
