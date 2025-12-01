package wallets

import (
	"context"
	"fmt"

	"tse-p3/db"

	"github.com/jackc/pgx/v5"
)

func GetUserWallets(ctx context.Context, userID int64) ([]WalletDescriptor, error) {
	var (
		query = `
			SELECT amount, symbol
			FROM wallets
			WHERE user_id = $1
			ORDER BY symbol`
		rows      pgx.Rows
		err       error
		wds       []WalletDescriptor
		amountStr string
		wd        WalletDescriptor
	)

	rows, err = db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&amountStr,
			&wd.Symbol,
		)
		if err != nil {
			return nil, err
		}

		_, err = fmt.Sscan(amountStr, &wd.Amount)
		if err != nil {
			return nil, err
		}

		wds = append(wds, wd)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return wds, nil
}

func CreateOrUpdateUserWallet(ctx context.Context, symbol string, amt string, userID int64) error {
	var (
		query = `
			INSERT INTO wallets (user_id, symbol, amount)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, symbol)
			DO UPDATE SET amount = EXCLUDED.amount`
		err error
	)

	_, err = db.Pool.Exec(ctx, query,
		userID,
		symbol,
		amt,
	)

	if err != nil {
		return err
	}

	return nil
}
