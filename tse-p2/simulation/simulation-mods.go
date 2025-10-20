package simulation

import (
    "fmt"
    "tse-p2/ledger"
    "tse-p2/candles"
    "tse-p2/exchange"
)

func (sim *Simulation) GetCandles() ([]candles.Candle, error) {
    sim.LedgerLock.Lock()
    defer sim.LedgerLock.Unlock()

    exItem, ok := sim.Ledger[sim.ExAddr]
    if !ok {
        return nil, fmt.Errorf("exchange not found on ledger")
    }
    ex, ok := exItem.(exchange.ConstantProductExchange)
    if !ok {
        return nil, fmt.Errorf("invalid exchange type")
    }

    auditItem, ok := sim.Ledger[ex.CndlAddr]
    if !ok {
        return nil, fmt.Errorf("candle audit not found on ledger")
    }
    audit, ok := auditItem.(candles.CandleAudit)
    if !ok {
        return nil, fmt.Errorf("invalid candle audit type")
    }

    candles := audit.CandleHistory.CandlesInOrder()
    candles = append(candles, audit.CurrentCandle)

    return candles, nil
}

func (s *Simulation) PlaceUserTrade(from, to string, confidence float64) error {
    var (
        tx  ledger.Tx
        err error
    )

    tx, err = s.CliTrader.CreateSwapTransaction(
        from,
        to,
        confidence,
    )
    if err != nil {
        return fmt.Errorf("failed to create tx: %v\n", err)
    }
    s.placeTx(tx)
    return nil
}


func (s *Simulation) placeTx(tx ledger.Tx) {
    (&s.MemoryPool).PushTx(tx)
}

func (s *Simulation) ledgerLookup(addr ledger.LedgerAddr) ledger.LedgerItem {
    return (*s).Ledger[addr]
}