package globals

import (
	"time"
	"github.com/holiman/uint256"
)

// ---- Token Settings ---- //
var TokenScaleFactor *uint256.Int = uint256.NewInt(1_000_000) // 10^6
var TokenScaleFactorf64 float64 = 1_000_000

// ----Default Currency Exchange ---- //
var USDSymbol = "usd"
var TSESymbol = "tse"
const USDCurrencyAmount = 1_000_000 	// Starting USD Exchange Amount
const TSECurrencyAmount = 1_000_000		// Starting TSE Exchange Amount
const UserStartingBalance = 10_000 		// Starting Balance Per User


// ----- Candle Audit Settings ------ // 
const DefaultAuditerBufferSize = 1000

// ------ MemoryPool Settings ----- //
const DefaultMemoryPoolSize = 200

// ------ Miner Settings ------- //
const MaxBlockSize = 7
const TimeBetweenBlocks = 2 * time.Millisecond

// ----- Bot Settings ----- //
const BotTaskDelay = 1 * time.Millisecond