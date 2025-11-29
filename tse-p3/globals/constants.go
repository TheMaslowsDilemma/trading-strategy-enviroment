package globals

import (
	"time"
	"github.com/holiman/uint256"
)

// ---- Token Settings ---- //
var TokenScaleFactor *uint256.Int = uint256.NewInt(1_000_000_000_000)
var TokenScaleFactorf64 float64 = 1_000_000_000_000

// ----Default Currency Exchange ---- //
var USDSymbol = "usd"
var TSESymbol = "tse"
const USDCurrencyAmount = 10_000_000_00	// Starting USD Exchange Amount
const TSECurrencyAmount = 100_000_000_000	// Starting TSE Exchange Amount
const UserStartingBalance = 1000 			// Starting Balance Per User


// ----- Candle Audit Settings ------ // 
const DefaultAuditerBufferSize = 500

// ------ MemoryPool Settings ----- //
const DefaultMemoryPoolSize = 200

// ------ Miner Settings ------- //
const MaxBlockSize = 21
const TimeBetweenBlocks = 250 * time.Millisecond

// ----- Bot Settings ----- //
const BotTaskDelay = 50 * time.Millisecond