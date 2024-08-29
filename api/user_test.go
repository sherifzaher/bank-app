package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mockdb "github.com/sherifzaher/clone-simplebank/db/mock"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type equalCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e equalCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.VerifyPassword(arg.HashedPassword, e.password)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e equalCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return equalCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       "Sherif Zaher",
		Email:          util.RandomEmail(),
	}

	return arg, password
}

func TestUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		buildStub     func(store *mockdb.MockStore)
		buildRequest  func() (string, string, io.Reader)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Create User",
			buildStub: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:       user.Username,
					HashedPassword: user.HashedPassword,
					Email:          user.Email,
					FullName:       user.FullName,
				}
				store.
					EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := fmt.Sprintf("/users")
				arg := CreateUserParams{
					Username: user.Username,
					Password: password,
					Email:    user.Email,
					FullName: user.FullName,
				}
				params, err := json.Marshal(arg)
				require.NoError(t, err)
				return url, http.MethodPost, bytes.NewReader(params)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user, password)
			},
		},
		{
			name: "Get User",
			buildStub: func(store *mockdb.MockStore) {
				arg := GetUserParams{
					Username: user.Username,
				}
				store.
					EXPECT().
					GetUser(gomock.Any(), arg.Username).
					Times(1).
					Return(user, nil)
			},
			buildRequest: func() (string, string, io.Reader) {
				url := fmt.Sprintf("/users/%v", user.Username)
				return url, http.MethodGet, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				//requireBodyMatchUser(t, recorder.Body, user, password)
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
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User, password string) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)

	validPassword := util.VerifyPassword(user.HashedPassword, password)
	require.NoError(t, validPassword)
}
