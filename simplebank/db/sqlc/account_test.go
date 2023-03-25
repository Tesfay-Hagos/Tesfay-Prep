package db

import (
	"context"
	"database/sql"
	"tesfayprep/simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateaccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currence: util.RandomCurrency(),
	}
	account, err := testQueries.Createaccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currence, account.Currence)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotZero(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currence, account2.Currence)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	accountupdated, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotZero(t, accountupdated)

	require.Equal(t, account.ID, accountupdated.ID)
	require.Equal(t, account.Owner, accountupdated.Owner)
	require.Equal(t, arg.Balance, accountupdated.Balance)
	require.Equal(t, account.Currence, accountupdated.Currence)
	//require.WithinDuration(t, account.CreatedAt, accountupdated.CreatedAt, time.Second)
}

func TestDelateAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account1, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account1)
}

func TestListAccounts(t *testing.T) {
	var lastaccount Account
	for k := 0; k < 10; k++ {
		lastaccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastaccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastaccount.Owner, account.Owner)
	}
}
