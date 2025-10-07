package simulation

import (
    "fmt"
    "tse-p2/ledger"
    "tse-p2/wallet"
    "tse-p2/token"
)

// -- Organize the Simulation Modification Logic -- //
func (s *Simulation) AddLedgerItem(id ledger.LedgerAddr, li ledger.LedgerItem) error {
    var existing ledger.LedgerItem

    s.LedgerLock.Lock()
    defer s.LedgerLock.Unlock()

    existing = s.Ledger[id]
    if existing != nil {
        return fmt.Errorf("ledger item already exists at %v", id)
    }

    s.Ledger[id] = li
    return nil
}


func (s *Simulation) AddWallet(initamnt uint64) ledger.LedgerAddr {
    var (
        walletAddr      ledger.LedgerAddr
        usdRsvAddr      ledger.LedgerAddr
        usdRsv          token.TokenReserve
        w               wallet.Wallet
        walletRsvs      []ledger.LedgerAddr
        err             error
    )
    
    for {
        // Add initial reserve to ledger
        usdRsvAddr = ledger.RandomLedgerAddr()
        usdRsv = token.TokenReserve {
            Amount: initamnt,
            Symbol: "usd",
        }
        err = s.AddLedgerItem(usdRsvAddr, usdRsv)
        if err != nil {
            continue
        }

        // Add the wallet to the ledger
        walletAddr = ledger.RandomLedgerAddr()
        walletRsvs = make([]ledger.LedgerAddr, 0)
        w = wallet.Wallet {
            TraderId: uint64(walletAddr),
            Reserves: walletRsvs,
        }
        w.AddReserve(usdRsvAddr)
        err = s.AddLedgerItem(walletAddr, w)
        if err != nil {
            continue
        }

        break
    }

    return walletAddr
}
