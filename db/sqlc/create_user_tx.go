package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreateUser func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (s *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := s.withTx(ctx, func(queries *Queries) error {
		var err error
		result.User, err = queries.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreateUser(result.User)
	})
	return result, err
}
