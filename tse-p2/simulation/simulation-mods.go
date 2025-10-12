package simulation

import (
    "fmt"
    "tse-p2/ledger"
    "tse-p2/wallet"
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


func (s *Simulation) AddWallet(amount float64) ledger.LedgerAddr {
    return wallet.InitWallet( "usd", amount, s.Ledger)
}
