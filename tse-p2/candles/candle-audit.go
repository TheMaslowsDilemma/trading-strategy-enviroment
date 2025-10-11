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

func NewCandleAudit(n int) CandleAudit {
    var cq CandleCirq = NewCandleCirq(n)
    return CandleAudit {
        LastActiveGrp: 0,
        Constructing: false,
        CurrentCandle: Candle{},
        PastCandle: cc,
    }
}

func (ca *CandleAudit) Add(uint64 tick, cost float64, volume uint64) {
    var (
        tmod    uint32
        tgrp      uint64
    )
    
    tgrp = tick / globals.TICK_PER_SECOND

    if tgrp != ca.LastActiveGrp {
        ca.startNextCandle(cost, volume)
    } else {
        &(ca.CurrentCandle).Add(cost, volume)
    }
    ca.LastActiveGrp = tgrp
}

func (ca *CandleAudit) startNextCandle(cost float64, volume uint64) {
    // NOTE when we enqueue a candle we loose older ones.
    // TODO give auditers a channel to emit candles to be stored
    // in some longer term storage like postgres db -- eventually
    
    if ca.Constructing {
        ca.CandleHistory.Enqueue(ca.CurrentCandle)
    }

    ca.CurrentCandle = Candle{}
    (&ca.CurrentCandle).Open(cost, volume)
    ca.Constructing = true
}

/** --- LedgerItem Interface Implementation --- **/

func (ca *CandleAudit) Copy() ledger.LedgerItem {
    return &CandleAudit {
        LastActive: ca.LastActive,
        Constructing: ca.Constructing,
        CurrentCandle: ca.CurrentCandle.Copy(),
        CandleHistory: ca.CandleHistory.Copy(),
    }
}


func (ca *CandleAudit) String() string {
    return fmt.Sprintf("{ cc: %v, last-active: %v }",
        ca.CurrentCandle.String(),
        ca.LastActive,
    )
}

func (ca *CandleAudit) Hash() [sha256.Size]byte {
    return sha256.Sum([]byte(ca.String()))
}
