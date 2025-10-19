package exchange

import (
    "fmt"
    "tse-p2/ledger"
    "tse-p2/wallet"
    "tse-p2/token"
    "tse-p2/candles"
)

type SwapExactTokensForTokensTx struct {
    SymbolIn            string
    SymbolOut           string
    AmountIn            float64
    AmountMinOut        float64
    WalletAddr          ledger.LedgerAddr
    ExchangeAddr        ledger.LedgerAddr
    Callback            func(result ledger.TxResult)
}

// Returns a Partial Ledger -- modifications only -- to be merged by the miner.
func (tx SwapExactTokensForTokensTx) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
    var (
        exchange                  *ConstantProductExchange
        wlt                       *wallet.Wallet
        sendReserveAddr           ledger.LedgerAddr
        recvReserveAddr           ledger.LedgerAddr
        sendReserve               *token.TokenReserve
        recvReserve               *token.TokenReserve
        exchangeReserveA          *token.TokenReserve
        exchangeReserveB          *token.TokenReserve
        exchangeCandleAudit       *candles.CandleAudit
        calculatedPrice           float64
        amountOut                 float64
        partialLedger             ledger.Ledger
        err                       error
    )

    wlt, err = wallet.WalletFromLedgerItem(l[tx.WalletAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet: %v", err)
    }
    
    sendReserveAddr, err = wlt.GetReserveAddr(tx.SymbolIn, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's input token reserve: %v", err)
    }
    sendReserve, err = token.TkrFromLedgerItem(l[sendReserveAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's input token reserve: %v", err)
    }

    recvReserveAddr, err = wlt.GetReserveAddr(tx.SymbolOut, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's output token reserve: %v", err)
    }
    recvReserve, err = token.TkrFromLedgerItem(l[recvReserveAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's output token reserve: %v", err)
    }

    exchange, err = CpeFromLedgerItem(l[tx.ExchangeAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange: %v", err) // Fixed: was missing %v
    }
    //exchange = (*exchange).Copy()

    exchangeReserveA, err = token.TkrFromLedgerItem(l[exchange.TkrAddrA])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's token reserve A: %v", err)
    }

    //exchangeReserveA = (*exchangeReserveA).Copy()

    exchangeReserveB, err = token.TkrFromLedgerItem(l[exchange.TkrAddrB])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's token reserve B: %v", err)
    }
    //exchangeReserveB = (*exchangeReserveB).Copy()
    
    exchangeCandleAudit, err = candles.CandleAuditFromLedgerItem(l[exchange.CndlAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's candle audit: %v", err)
    }

    if sendReserve.Amount < tx.AmountIn {
        return nil, fmt.Errorf("insufficient funds for swap")
    }
    
    // Find out how we want to swap & if we even can
    
    partialLedger = make(ledger.Ledger)
    if exchangeReserveA.Symbol == tx.SymbolIn && exchangeReserveB.Symbol == tx.SymbolOut {
        amountOut, err = exchange.SwapAForB(l, tx.AmountIn)
        if err != nil {
            return nil, fmt.Errorf("swap A for B failed: %v", err)
        }

        if amountOut < tx.AmountMinOut {
            return nil, fmt.Errorf("swap slippage too high: got %f, minimum required %f", amountOut, tx.AmountMinOut)
        }

        if amountOut > exchangeReserveB.Amount {
            return nil, fmt.Errorf("swap output greater then exchange B reserve, %v > %v", amountOut, exchangeReserveB.Amount)
        }
       
        //fmt.Printf("applying swap: %v %v -> %v %v\n", tx.SymbolIn, tx.AmountIn, tx.SymbolOut, amountOut)

        // Transfer input tokens: wallet -> exchange
        sendReserve.Amount -= tx.AmountIn
        exchangeReserveA.Amount += tx.AmountIn
        
        // Transfer output tokens: exchange -> wallet
        exchangeReserveB.Amount -= amountOut
        recvReserve.Amount += amountOut
    } else if exchangeReserveB.Symbol == tx.SymbolIn && exchangeReserveA.Symbol == tx.SymbolOut {
        amountOut, err = exchange.SwapBForA(l, tx.AmountIn)
        if err != nil {
            return nil, fmt.Errorf("swap B for A failed: %v", err)
        }

        if amountOut < tx.AmountMinOut {
            return nil, fmt.Errorf("swap slippage too high: got %f, minimum required %f", amountOut, tx.AmountMinOut)
        }

        if amountOut > exchangeReserveA.Amount {
            return nil, fmt.Errorf("swap output greater then exchange A reserve, %v > %v", amountOut, exchangeReserveA.Amount)
        }

        //fmt.Printf("applying swap: %v %v -> %v %v\n", tx.SymbolIn, tx.AmountIn, tx.SymbolOut, amountOut)

        // Transfer input tokens: wallet -> exchange  
        sendReserve.Amount -= tx.AmountIn
        exchangeReserveB.Amount += tx.AmountIn // Input goes to reserve B
        
        // Transfer output tokens: exchange -> wallet
        exchangeReserveA.Amount -= amountOut // Output comes from reserve A
        recvReserve.Amount += amountOut
    } else {
        return nil, fmt.Errorf("symbol mismatch: tx{%s -> %s} != exchange{%s <-> %s}",
                tx.SymbolIn, tx.SymbolOut, exchangeReserveA.Symbol, exchangeReserveB.Symbol)
    }

    // Finally Update the Partial Ledger & Calculate new spot price
    partialLedger[sendReserveAddr] = *sendReserve
    partialLedger[recvReserveAddr] = *recvReserve
    partialLedger[exchange.TkrAddrA] = *exchangeReserveA
    partialLedger[exchange.TkrAddrB] = *exchangeReserveB
    
    calculatedPrice, err = exchange.GetPriceB(partialLedger)
    if err != nil {
        return nil, err
    } 
    exchangeCandleAudit.Add(tick, calculatedPrice, tx.AmountIn)

    partialLedger[exchange.CndlAddr] = exchangeCandleAudit
    return partialLedger, nil
}


func (tx SwapExactTokensForTokensTx) Notify(result ledger.TxResult) {
    tx.Callback(result)
}