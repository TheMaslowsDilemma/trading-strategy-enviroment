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
}

// Returns a Partial Ledger -- modifications only -- to be merged by the miner.
func (tx SwapExactTokensForTokensTx) Apply(tick uint64, l ledger.Ledger) (ledger.Ledger, error) {
    var (
        exg             *ConstantProductExchange
        w               *wallet.Wallet
        waddrO          ledger.LedgerAddr   // wallet tk "out" reserveaddr
        waddrI          ledger.LedgerAddr   // wallet tk "in"  reserveaddr
        wtkrO           *token.TokenReserve
        wtkrI           *token.TokenReserve
        etkrA           *token.TokenReserve // exchange A reserve
        etkrB           *token.TokenReserve // exchange B reserve
        auditer         *candles.CandleAudit // exchange audit
        price            float64
        amtO            float64
        ldgp            ledger.Ledger
        err             error
    )

    w, err = wallet.WalletFromLedgerItem(l[tx.WalletAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet: %v", err)
    }
    
    waddrI, err = w.GetReserveAddr(tx.SymbolIn, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's source reserve: %v", err)
    }
    wtkrI, err = token.TkrFromLedgerItem(l[waddrI])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's source reserve %v", err)
    }

    waddrO, err = w.GetReserveAddr(tx.SymbolOut, l)
    if err != nil {
        return nil, fmt.Errorf("failed to find wallet's destination reserve: %v", err)
    }
    wtkrO, err = token.TkrFromLedgerItem(l[waddrO])
    if err != nil {
        return nil, fmt.Errorf("failed to cast wallet's destination reserve: %v", err)
    }

    exg, err = CpeFromLedgerItem(l[tx.ExchangeAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange: ", err)
    }

    etkrA, err = token.TkrFromLedgerItem(l[exg.TkrAddrA])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's tkr-a: %v", err)
    }

    etkrB, err = token.TkrFromLedgerItem(l[exg.TkrAddrB])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's tkr-b: %v", err)
    }
    
    auditer, err = candles.CandleAuditFromLedgerItem(l[exg.CndlAddr])
    if err != nil {
        return nil, fmt.Errorf("failed to cast exchange's auditer: %v", err)
    }

    if wtkrI.Amount < tx.AmountIn {
        return nil, fmt.Errorf("insufficient funds for swap")
    }
    
    // Find out how we want to swap & if we even can
    if etkrA.Symbol == tx.SymbolIn && etkrB.Symbol == tx.SymbolOut {
        amtO, err = exg.SwapAForB(l, tx.AmountIn)
        if err != nil {
            return nil, fmt.Errorf("swap a for b failed: %v", err)
        }

        if amtO < tx.AmountMinOut {
            return nil, fmt.Errorf("swap slippage too high")
        }
        
        // Make the Diff Ledger
        ldgp = make(ledger.Ledger)

        // Calculate Price & Audit -- A with respect to B
        price = amtO / tx.AmountIn
        auditer.Add(tick, price, tx.AmountIn) // Volume on an exchange is defined by volume of A

        // Move In Funds to the Exchange
        wtkrI.Amount -= tx.AmountIn
        etkrA.Amount += tx.AmountIn
        
        // Move Out Funds to the Wallet
        etkrB.Amount -= amtO
        wtkrO.Amount += amtO

        // Create Diff Ledger & Return
        ldgp[waddrI] = wtkrI
        ldgp[waddrO] = wtkrO
        ldgp[exg.TkrAddrA] = etkrA
        ldgp[exg.TkrAddrB] = etkrB
        ldgp[exg.CndlAddr] = auditer // TODO copy?
        
        return ldgp, nil
        
    } else if etkrB.Symbol == tx.SymbolIn && etkrA.Symbol == tx.SymbolOut {
        amtO, err = exg.SwapBForA(l, tx.AmountIn)
        if err != nil {
            return nil, fmt.Errorf("swap b for a failed: %v", err)
        }

        if amtO < tx.AmountMinOut {
            return nil, fmt.Errorf("swap slippage too high")
        }

        // Make the Diff Ledger
        ldgp = make(ledger.Ledger)

        // Calculate Price & Audit -- A with respect to B
        price = tx.AmountIn / amtO // calcuting A with respec to B
        auditer.Add(tick, price, amtO) // Volume on an exchange is defined by volume of A

        // Move In Funds to the Exchange
        wtkrI.Amount -= tx.AmountIn
        etkrA.Amount += tx.AmountIn
        
        // Move Out Funds to the Wallet
        etkrB.Amount -= amtO
        wtkrO.Amount += amtO

        // Create Diff Ledger & Return
        ldgp[waddrI] = wtkrI
        ldgp[waddrO] = wtkrO
        ldgp[exg.TkrAddrA] = etkrA
        ldgp[exg.TkrAddrB] = etkrB
        ldgp[exg.CndlAddr] = auditer // TODO copy?

        return ldgp, nil    
    }
    return nil, fmt.Errorf("failed to match symbols tx{ %v -> %v } != ex{ %v <-> %v }",
        tx.SymbolIn,
        tx.SymbolOut,
        etkrA.Symbol,
        etkrB.Symbol,
    )
}
