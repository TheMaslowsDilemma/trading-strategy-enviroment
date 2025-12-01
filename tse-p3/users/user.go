package users

import (
	"fmt"
	"errors"
	"strconv"
	"context"

	"tse-p3/db"
	"tse-p3/globals"
	"tse-p3/ledger"
	"tse-p3/traders"
	"tse-p3/wallets"
	"tse-p3/simulation"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID				int64	`json:"id"`
	Name			string	`json:"name"`
	TraderID		uint64	`json:"trader_id"`
	PasswordHash	string	`json:"-"`
	Active			bool	`json:"active"`
}

func CreateUser(ctx context.Context, username, password string, sim *simulation.Simulation) error {
	var (
		hash		[]byte
		trdr_id		uint64
		query		string
		userID		int64
		err			error
	)
		
	trdr_id = globals.Rand64()

	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO users (name, password_hash, trader_id, active)
		VALUES ($1, $2, $3, false)
		RETURNING id`

	err = db.Pool.QueryRow(ctx, query, username, string(hash), strconv.FormatUint(trdr_id, 10)).Scan(&userID)
	if err != nil {
		return err
	}

	err = wallets.CreateOrUpdateUserWallet(ctx, globals.USDSymbol, strconv.FormatUint(globals.UserStartingBalance, 10), userID)
	if err != nil {
		return err
	}
	return nil
}

func (user User) BeginUserSession(ctx context.Context, s *simulation.Simulation) error {
	var (
		trader		*traders.Trader
		wlt_addr	ledger.Addr
		wds			[]wallets.WalletDescriptor
		wd			wallets.WalletDescriptor
		err			error
	)

	wds, err = wallets.GetUserWallets(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get user wallet: %v", err)
	}

	trader = traders.CreateTraderWithId(user.Name, user.TraderID)
	s.AddTrader(trader)

	for _, wd = range wds {
		wlt_name := fmt.Sprintf("%v:w:%v", user.Name, wd.Symbol)
		wd.Name = wlt_name
		wlt_addr = s.AddWallet(wd)
		trader.AddWallet(wd.Symbol, wlt_addr)
	}

	return nil
}

// THIS could be re thought... were locking up the sim while doing this...
func (user User) CleanUserSession(ctx context.Context, s *simulation.Simulation) {
	var (
		trader		*traders.Trader
		wlt_addr	ledger.Addr
		wlt			wallets.Wallet
		err			error
	)
	trader = s.Traders[user.TraderID]
	if trader == nil {
		return
	}

	s.PrimaryLock.Lock()
	s.SecondaryLock.Lock()
	for _, wlt_addr = range trader.Wallets {
		wlt = s.PrimaryLedger.GetWallet(wlt_addr)
		err = wallets.CreateOrUpdateUserWallet(
			ctx, 
			wlt.Reserve.Symbol,
			wlt.Reserve.Amount.String(), 
			user.ID,
		)
		if err != nil {
			fmt.Printf("failed to update wallet for '%v': %v\n", user.Name, err)
		}
		s.PrimaryLedger.RemoveWallet(wlt_addr)
		s.SecondaryLedger.RemoveWallet(wlt_addr)
		s.PrimaryLedger.RemoveWallet(wlt_addr)
	}

	s.PrimaryLock.Unlock()
	s.SecondaryLock.Unlock()

	delete(s.Traders, user.TraderID)
}

func GetUserByName(ctx context.Context, name string) (User, error) {
	var (
		u			User
		query		string
		tid_str		string
		err			error
	)

	query = `
		SELECT id, name, trader_id, password_hash, active
		FROM users
		WHERE name = $1`


	err = db.Pool.QueryRow(ctx, query, name).Scan(
		&u.ID,
		&u.Name,
		&tid_str,
		&u.PasswordHash,
		&u.Active,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return u, errors.New("user not found")
	}

	if err != nil {
		return u, err
	}

	u.TraderID, err = strconv.ParseUint(tid_str, 10, 64)
	if err != nil {
		return u, err
	}
	
	return u, nil
}

func GetUserById(ctx context.Context, id int64) (User, error) {
	var (
		u			User
		query		string
		tid_str		string
		err			error
	)

	query = `
		SELECT id, name, trader_id, password_hash, active
		FROM users
		WHERE id = $1`

	err = db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&tid_str,
		&u.PasswordHash,
		&u.Active,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return u, errors.New("user not found")
	}

	if err != nil {
		return u, err
	}

	u.TraderID, err = strconv.ParseUint(tid_str, 10, 64)
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