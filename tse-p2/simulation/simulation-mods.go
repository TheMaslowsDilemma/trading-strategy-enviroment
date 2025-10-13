package simulation

import (
    "fmt"
    "tse-p2/ledger"
)
func (s *Simulation) PlaceUserTrade(from, to string, confidence float64) {
    var (
        tx  ledger.Tx
        err error
    )

    tx, err = s.CliTrader.CreateSwapTransaction(
        from,
        to,
        confidence,
        s.Ledger,
    )
    if err != nil {
        fmt.Printf("failed to create tx: %v\n", err)
        return
    }
    (&s.MemoryPool).PushTx(tx)
}
