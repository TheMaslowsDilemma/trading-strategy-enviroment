package trader

import (
    "tse-p2/ledger"
    "tse-p2/candles"
    "tse-p2/strategy"
)

type Trader struct {
    Wallet      ledger.Addr
    Decider     strategy.Strategy
    PendingTx   chan uint8
    Logs        chan string
}

func CreateTrader(s strategy.Strategy, logsize int) Trader {
    var (
        ptx     chan uint8
        lgs     chan string
        waddr   ledger.Addr
    )

    waddr = ledger.RandomLedgerAddr()
    ptx = make(chan uint8, 1)
    lgs = make(chan string, logsize)

    return Trader {
        Wallet: waddr,
        Decider: s,
        PendingTx: ptx,
        Logs: lgs,
    }
}

func (t *Trader) MakeDecision(sym string, cs []candle.Candle, l ledger.Ledger) *ledger.Tx {

    select {
    case <-t.PendingTx:
        return nil
    default:
        return t.getTx(sym, cs, l)
    }
}

func (t *Trader) getTx(cs []candles.Candle, l ledger.Ledger) (*ledger.Tx, error) {
    var (
        action strategy.Action
        confidence float64
        err error
    )

    action, confidence = t.Decider.Decide(cs, l)

    switch action {
    case strategy.Hold:
        return nil, nil
    case strategy.Buy:
        return t.buyTx(sym, confidence, l)
    case strategy.Sell:
        return t.sellTx(sym, confidence, l)
    }
}

func (t *Trader) buyTx(sym string, c float64, l *ledger.Ledger) (*ledger.Tx, error) {
    var (
        wlt     *wallet.Wallet
        tkaddrA ledger.LedgerAddr
        tkaddrB ledger.LedgerAddr
        tkrA    *token.TokenReserve
        tkrB    *token.TokenReserve
        err error
    )
    
    wlt, err = wallet.WalletFromLedgerItem(l[t.Wallet])
    if err != nil {
         return nil, fmt.Errorf("buytx failed to cast wallet: %v", err)
    }

    tkAddr, err = wlt.GetReserveAddr(sym, l)
    if err != nil {
        return ni, fmt.Errorf("buytx failed to get token reserve: %v", err)
    }
    
    tkA, err = token.TkrFromLedgerItem(l[tkAddr])
    if err != nil {
        return nil, fmt.Errorf("buytx failed to cast tkA: %v", err)
    }
    
     
}
