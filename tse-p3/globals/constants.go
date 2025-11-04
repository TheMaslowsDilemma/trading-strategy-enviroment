package globals

import (
	"time"
)

// ------- Default Currency Exchange --------- //
var USDCurrencySymbol = "usd"
const USDCurrencyAmount = 1_000_000_000
var TSECurrencySymbol = "tsd"
const TSECurrencyAmount = 1_000_000_000

const UserStartingBalance = 10_000


// ----- Candle Audit Settings ------ // 
const DefaultAuditerBufferSize = 1000

// ------ MemoryPool Settings ----- //
const DefaultMemoryPoolSize = 200

// ------ Miner Settings ------- //
const MaxBlockSize = 7
const TimeBetweenBlocks = 100 * time.Millisecond