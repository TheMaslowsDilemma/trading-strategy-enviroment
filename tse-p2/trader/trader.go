package trader

import (
    "tse-p2/ledger"
    "tse-p2/candles"
    "tse-p2/strategy"
    "tse-p2/wallet"
    "tse-p2/exchange"
)

type Trader struct {
    ExAddr      ledger.Addr
    WalletAddr  ledger.Addr
    Decider     strategy.Strategy
    PendingTx   bool
    Logs        chan string
}

/** NOTE for now, traders are associated with a specific exchange **/
func InitTrader(s strategy.Strategy, logsize int, eaddr ledger.Addr, l ledger.Ledger) Trader {
    var (
        ptx     bool
        lgs     chan string
        waddr   ledger.Addr
        wlt     wallet.Wallet
    )
    ptx = false
    lgs = make(chan string, logsize)
    waddr = wallet.InitWallet("usd", 10000, l) // NOTE default seed amount 10000

    return Trader {
        ExAddr: eaddr
        WalletAddr: waddr,
        Decider: s,
        PendingTx: ptx,
        Logs: lgs,
    }
}

func (t *Trader) MakeDecision(sym string, cs []candle.Candle, l ledger.Ledger) *ledger.Tx {
    var (
        tx  *ledger.Tx
        err error
    )

    if t.PendingTx {
        return nil
    }

    tx, err = t.getTx(sym, cs, l)
    if err != nil {
        // TODO ? log err
        return nil
    }

    return tx
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

func (t *Trader) buyTx(sym string, confidence float64, l *ledger.Ledger) (*ledger.Tx, error) {
    var (
        wlt         *wallet.Wallet
        tkaddrA     ledger.LedgerAddr
        tkaddrB     ledger.LedgerAddr
        tkrA        *token.TokenReserve
        tkrB        *token.TokenReserve
        amtIn       float64
        amtMinOut   float64
        err         error
    )
    
    wlt, err = wallet.WalletFromLedgerItem(l[t.WalletAddr])
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
    
    amtIn  = tkA.Amount * confidence
    amtMinOut = amntIn * 

    return &SwapExactTokensForTokensTx {
        SymbolIn: "usd", // whenever we buy we use "usd"
        SymbolOut: sym,
        AmountIn: amtIn,
        AmountMinOut: amtOut,
        WalletAddr: t.Wallet,
        ExchangeAddr: t.ExAddr,
    }
}
