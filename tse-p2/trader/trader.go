package trader

import (
    "fmt"
    "tse-p2/token"
    "tse-p2/ledger"
    "tse-p2/candles"
    "tse-p2/strategy"
    "tse-p2/wallet"
    "tse-p2/exchange"
)

type Trader struct {
    ExAddr      ledger.LedgerAddr
    ExSymA      string
    ExSymB      string
    WalletAddr  ledger.LedgerAddr
    Decider     strategy.Strategy
    PendingTx   bool
    Logs        chan string
}

/** NOTE for now, traders are associated with a specific exchange **/
func CreateTrader(s strategy.Strategy, ls int, waddr, eaddr ledger.LedgerAddr, symA, symB string, l ledger.Ledger) Trader {
    var (
        ptx     bool
        lgs     chan string
    )
    ptx = false
    lgs = make(chan string, ls)

    return Trader {
        ExAddr: eaddr,
        ExSymA: symA,
        ExSymB: symB,
        WalletAddr: waddr,
        Decider: s,
        PendingTx: ptx,
        Logs: lgs,
    }
}

func (t *Trader) MakeDecision(cs []candles.Candle, l ledger.Ledger) ledger.Tx {
    var (
        tx  ledger.Tx
        err error
    )

    if t.PendingTx {
        return nil
    }

    tx, err = t.getTx(cs, l)
    if err != nil {
        // TODO ? log err
        return nil
    }

    return tx
}

func (t *Trader) getTx(cs []candles.Candle, l ledger.Ledger) (ledger.Tx, error) {
    var (
        action strategy.Action
        confidence float64
    )

    action, confidence = t.Decider.Decide(cs, l)
 
    switch action {
    case strategy.Hold:
        return nil, nil
    case strategy.Buy:
        return t.SwapTx(t.ExSymA, t.ExSymB, confidence, l)
    case strategy.Sell:
        return t.SwapTx(t.ExSymB, t.ExSymA, confidence, l)
    }
    return nil, nil
}

func (t *Trader) SwapTx(sndSym, rcvSym string, confidence float64, l ledger.Ledger) (ledger.Tx, error) {
    var (
        wlt         *wallet.Wallet
        sndAddr     ledger.LedgerAddr
        sndTkr      *token.TokenReserve
        amtIn       float64
        amtOutMin   float64
        err         error
    )
    
    wlt, err = wallet.WalletFromLedgerItem(l[t.WalletAddr])
    if err != nil {
         return nil, fmt.Errorf("buytx failed to cast wallet: %v", err)
    }

    sndAddr, err = wlt.GetReserveAddr(sndSym, l)
    if err != nil {
        return nil, fmt.Errorf("buytx failed to get token reserve: %v", err)
    }
    
    sndTkr, err = token.TkrFromLedgerItem(l[sndAddr])
    if err != nil {
        return nil, fmt.Errorf("buytx failed to cast sendTkr: %v", err)
    }

    amtIn = sndTkr.Amount * confidence
    amtOutMin = 0 // NOTE fix this in the future to use descrim TODO fix

    return exchange.SwapExactTokensForTokensTx {
        SymbolIn: sndSym,
        SymbolOut: rcvSym,
        AmountIn: amtIn,
        AmountMinOut: amtOutMin,
        WalletAddr: t.WalletAddr,
        ExchangeAddr: t.ExAddr,
    }, nil
}
