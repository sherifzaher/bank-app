package api

import (
	"bytes"
	"database/sql"
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

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore)
		buildRequest  func() (string, string, io.Reader)
		authUsername  string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK - code=200",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
						Balance:  0,
					}).
					Times(1).
					Return(account, nil)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := "/accounts"
				arg := CreateAccountParams{
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
		{
			name:         "MissedParams - code=400",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := "/accounts"
				//arg := CreateAccountParams{
				//	Owner:    account.Owner,
				//	Currency: account.Currency,
				//}
				//params, err := json.Marshal(arg)
				//require.NoError(t, err)
				return url, http.MethodPost, bytes.NewReader(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:         "Invalid Currency - code=400",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := "/accounts"
				arg := CreateAccountParams{
					Currency: "USDS",
				}
				params, err := json.Marshal(arg)
				require.NoError(t, err)
				return url, http.MethodPost, bytes.NewReader(params)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			tc.buildStub(store)

			server := newTestStore(t, store)

			url, method, body := tc.buildRequest()
			request, err := http.NewRequest(method, url, body)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			createAndSetAuthToken(t, request, server.tokenMaker, tc.authUsername)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListAccountsAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts := make([]db.Account, n)

	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	type Query struct {
		PageID   int32
		PageSize int32
	}

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore)
		query         Query
		authUsername  string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  5,
					Offset: 0,
				}
				store.
					EXPECT().
					ListAccounts(gomock.Any(), arg).
					Times(1).
					Return(accounts, nil)
			},
			query: Query{PageID: 1, PageSize: int32(n)},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:         "Internal Server Error",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			query: Query{PageID: 1, PageSize: int32(n)},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:         "Invalid Page Size",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			query: Query{PageID: 1, PageSize: 60000},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "Invalid Page Number",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			query: Query{PageID: -1, PageSize: int32(n)},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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
			server := newTestStore(t, store)
			tc.buildStub(store)

			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			q.Add("page_number", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("owner", user.Username)
			request.URL.RawQuery = q.Encode()

			createAndSetAuthToken(t, request, server.tokenMaker, tc.authUsername)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore)
		authUsername  string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			authUsername: user.Username,
			buildStub: func(store *mockdb.MockStore) {
				arg := db.GetAccountParams{
					Owner:    user.Username,
					ID:       account.ID,
					Currency: account.Currency,
				}
				store.
					EXPECT().
					GetAccount(gomock.Any(), arg).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		//{
		//	name: "Invalid Page Number",
		//	buildStub: func(store *mockdb.MockStore) {
		//		store.
		//			EXPECT().
		//			ListAccounts(gomock.Any(), gomock.Any()).
		//			Times(0)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusBadRequest, recorder.Code)
		//	},
		//},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestStore(t, store)
			tc.buildStub(store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%v", account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("owner", user.Username)
			q.Add("currency", account.Currency)
			request.URL.RawQuery = q.Encode()

			createAndSetAuthToken(t, request, server.tokenMaker, tc.authUsername)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
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
