package users

import (
	"context"
	"errors"

	"tse-p3/db"
	"tse-p3/globals"
	"tse-p3/ledger"
	"tse-p3/simulation"
	"tse-p3/traders"
	"tse-p3/wallets"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type DataSubscription struct {
	Name  string            `json:"name"`
	Etype ledger.EntityType `json:"etype"`
	Addr  ledger.Addr       `json:"addr"`
}

type User struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	TraderID          uint64             `json:"trader_id"`
	PasswordHash      string             `json:"-"` // this wont be marshalled
	DataSubscriptions []DataSubscription `json:"data_subscriptions"`
}

func CreateUser(ctx context.Context, username, password string, sim *simulation.Simulation) (int64, error) {
	var (
		hash     []byte
		err      error
		trader   *traders.Trader
		wltDesc  wallets.WalletDescriptor
		wltAddr  ledger.Addr
		query    string
		userID   int64
	)

	trader = traders.CreateTrader()

	wltDesc = wallets.WalletDescriptor{
		Amount: globals.UserStartingBalance,
		Symbol: globals.USDSymbol,
	}
	wltAddr = sim.AddWallet(wltDesc)
	trader.AddWallet(wltDesc.Symbol, wltAddr)
	sim.AddTrader(trader)
	sim.AddWallet(wltDesc, fmt.Sprintf("%v:w:%v", username, globals.USDSymbol)

	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	query = `
		INSERT INTO users (name, password_hash, trader_id, data_subscriptions)
		VALUES ($1, $2, $3, '{}')
		RETURNING id`

	err = db.Pool.QueryRow(ctx, query, username, string(hash), trader.ID).Scan(&userID)
	if err != nil {
		return -1, err
	}

	ctx = context.WithValue(ctx, "user.id", userID)
	ctx = context.WithValue(ctx, "user.name", username)
	ctx = context.WithValue(ctx, "user.trader.id", trader.ID)

	return  nil
}

func GetUserByName(ctx context.Context, name string) (User, error) {
	var (
		u     User
		query string
		err   error
	)

	query = `
		SELECT id, name, trader_id, password_hash, data_subscriptions
		FROM users
		WHERE name = $1`

	err = db.Pool.QueryRow(ctx, query, name).Scan(
		&u.ID,
		&u.Name,
		&u.TraderID,
		&u.PasswordHash,
		(*pgx.Jsonb)(&u.DataSubscriptions),
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
		u     User
		query string
		err   error
	)

	query = `
		SELECT id, name, trader_id, password_hash, data_subscriptions
		FROM users
		WHERE id = $1`

	err = db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.TraderID,
		&u.PasswordHash,
		(*pgx.Jsonb)(&u.DataSubscriptions),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return u, errors.New("user not found")
	}
	if err != nil {
		return u, err
	}

	return u, nil
}

func (u *User) ComparePassword(password string) bool {
	var err error
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) UpdateSubscriptions(ctx context.Context) error {
	var (
		query string
		err   error
	)

	query = `UPDATE users SET data_subscriptions = $1 WHERE id = $2`
	_, err = db.Pool.Exec(ctx, query, u.DataSubscriptions, u.ID)
	return err
}