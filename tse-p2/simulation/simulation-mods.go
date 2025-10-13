package simulation

import (
    "fmt"
    "tse-p2/ledger"
)
func (s *Simulation) PlaceUserTrade(from, to string, confidence float64) error {
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
        return fmt.Errorf("failed to create tx: %v\n", err)
    }
    (&s.MemoryPool).PushTx(tx)
    return nil
}
