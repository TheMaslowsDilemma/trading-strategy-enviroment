package simulation

import (
    "tse-p2/ledger"
)

func (s *Simulation) GetLedgerItemString(id ledger.LedgerAddr) (string, error) {
    var (
        str string
        err error
    )

    s.LedgerLock.Lock()
    str, err = s.Ledger.GetItemString(id)
    s.LedgerLock.Unlock()

    return str, err
}
