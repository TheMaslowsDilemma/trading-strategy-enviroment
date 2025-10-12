package candles

import (
    "fmt"
    "crypto/sha256"
    "tse-p2/globals"
    "tse-p2/ledger"
)

type CandleAudit struct {
    LastActive          uint64
    Constructing        bool
    CurrentCandle       Candle
    CandleHistory       CandleCirq
}

func InitCandleAudit(size uint32, l ledger.Ledger) ledger.LedgerAddr {
    var laddr ledger.LedgerAddr

    laddr = ledger.RandomLedgerAddr()
    l[laddr] = CandleAudit {
        LastActive: 0,
        Constructing: false,
        CurrentCandle: Candle{},
        CandleHistory: NewCandleCirq(size),
    }
    return laddr
}

func (ca *CandleAudit) Add(tick uint64, price float64, volume float64) {
    var timegroup uint64
    
    timegroup = tick / globals.TICK_PER_SECOND

    if timegroup != ca.LastActive {
        ca.startNextCandle(price, volume)
    } else {
        (&ca.CurrentCandle).Add(price, volume)
    }

    ca.LastActive = timegroup
}

func (ca *CandleAudit) startNextCandle(price float64, volume float64) {
    // NOTE when we enqueue a candle we loose older ones.
    // TODO give auditers a channel to emit candles to be stored
    // in some longer term storage like postgres db -- eventually
    
    if ca.Constructing {
        ca.CandleHistory.Enqueue(ca.CurrentCandle)
    }

    ca.CurrentCandle = Candle{}
    (&ca.CurrentCandle).Start(price, volume)
    ca.Constructing = true
}

/*** --- LedgerItem Implementation --- ***/

func (ca CandleAudit) Copy() ledger.LedgerItem {
    return CandleAudit {
        LastActive: ca.LastActive,
        Constructing: ca.Constructing,
        CurrentCandle: ca.CurrentCandle.Copy(),
        CandleHistory: ca.CandleHistory.Copy(),
    }
}

func (ca CandleAudit) String() string {
    return fmt.Sprintf("{ crnt-cndl: %v, last-active: %v }",
        ca.CurrentCandle.String(),
        ca.LastActive,
    )
}

func (ca CandleAudit) Hash() [sha256.Size]byte {
    return sha256.Sum256([]byte(ca.String()))
}

func CandleAuditFromLedgerItem(li ledger.LedgerItem) (*CandleAudit, error) {
    var (
        ca  CandleAudit
        ok  bool
    )

    if li == nil {
        return nil, fmt.Errorf("cannot cast candle audit from nil ledger item.")
    }
    
    if ca, ok = li.(CandleAudit); ok {
        return &ca, nil
    }
    return nil, fmt.Errorf("cannot cast candle audit from non candle audit ledger item.")
}
