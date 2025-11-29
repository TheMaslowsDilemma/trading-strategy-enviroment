package users

import (
	"fmt"
	"errors"
	"context"

	"tse-p3/db"
	"tse-p3/globals"
	"tse-p3/ledger"
	"tse-p3/simulation"
	"tse-p3/traders"
	"tse-p3/wallets"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID				int64	`json:"id"`
	Name			string	`json:"name"`
	TraderID		uint64	`json:"trader_id"`
	PasswordHash	string	`json:"-"` // this wont be marshalled
	Active			bool	`json:"active"`
}

func CreateUser(ctx context.Context, username, password string, sim *simulation.Simulation) error {
	var (
		hash		[]byte
		trader		*traders.Trader
		wlt_dsc		wallets.WalletDescriptor
		wlt_addr	ledger.Addr
		query		string
		userID		int64
		err			error
	)

	trader = traders.CreateTrader(username)

	wlt_dsc = wallets.WalletDescriptor{
		Name:	fmt.Sprintf("%v:w:%v", username, globals.USDSymbol),
		Amount:	globals.UserStartingBalance,
		Symbol:	globals.USDSymbol,
	}

	wlt_addr = sim.AddWallet(wlt_dsc)
	trader.AddWallet(wlt_dsc.Symbol, wlt_addr)
	sim.AddTrader(trader)


	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users (name, password_hash, trader_id)
		VALUES ($1, $2, $3)
		RETURNING id`

	err = db.Pool.QueryRow(ctx, query, username, string(hash), trader.Id).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByName(ctx context.Context, name string) (User, error) {
	var (
		u			User
		query		string
		err			error
	)

	query = `
		SELECT id, name, trader_id, password_hash, active
		FROM users
		WHERE name = $1`

	err = db.Pool.QueryRow(ctx, query, name).Scan(
		&u.ID,
		&u.Name,
		&u.TraderID,
		&u.PasswordHash,
		&u.Active,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return u, errors.New("user not found")
	}

	if err != nil {
		return u, err
	}

	return u, nil
}

func GetUserById(ctx context.Context, id int64) (User, error) {
	var (
		u			User
		query		string
		err			error
	)

	query = `
		SELECT id, name, trader_id, password_hash, active
		FROM users
		WHERE id = $1`

	err = db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.TraderID,
		&u.PasswordHash,
		&u.Active,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return u, errors.New("user not found")
	}

	if err != nil {
		return u, err
	}

	return u, nil
}

func SetUserActivity(ctx context.Context, uid int64, active bool) error {
	var (
		query	string
		err		error
	)
	query = `
		UPDATE users
		SET active = $1
		WHERE id = $2`

	_, err = db.Pool.Exec(ctx, query, active, uid)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	var err error
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}