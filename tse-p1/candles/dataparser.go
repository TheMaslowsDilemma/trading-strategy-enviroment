package candles

import (
	"time"
	"encoding/csv"
	"strconv"
	"os"
)

type CandleOrder uint8
const (
	OldestFirst CandleOrder = iota
	OldestLast
)

type CandlesCsvDescriptor struct {
	TimestampIndex int
	OpenIndex int
	HighIndex int
	LowIndex  int
	CloseIndex int
	VolumeIndex int
	ColumnCount int
	Order CandleOrder
}

func panicCheck(e error) {
    if e != nil {
        panic(e)
    }
}

// TODO: consider logging failed rows etc.
func LoadCandlesFromCsv(fp string, dsc CandlesCsvDescriptor) ([]Candle, error) {
	var (
		f *os.File
		rdr *csv.Reader
		rows [][]string
		rowcount int
		i int
		mi int
		cdls []Candle
		err error
	)

	f, err = os.Open(fp)
	panicCheck(err)
	defer f.Close()

	rdr = csv.NewReader(f)
	rows, err = rdr.ReadAll() // TODO: consider reading in batches
	panicCheck(err)

	rowcount = len(rows)

	mi = 1
	for mi < (rowcount - 1) {
		var (
			row []string
			ts time.Time
			open float64
			high float64
			low float64
			close float64
			volume float64
		)

		if dsc.Order == OldestFirst {
			i = mi
		} else {
			i = rowcount - mi
		}
		mi = mi + 1
		row = rows[i]

		if len(row) != dsc.ColumnCount {
			continue
		}

		ts, err := time.Parse("2006-01-02 15:04:05", row[dsc.TimestampIndex])
		if err != nil {
			continue
		}

		open, err = strconv.ParseFloat(row[dsc.OpenIndex], 64)
		if err != nil {
			continue
		}

		close, err = strconv.ParseFloat(row[dsc.CloseIndex], 64)
		if err != nil {
			continue
		}

		high, err = strconv.ParseFloat(row[dsc.HighIndex], 64)
		if err != nil {
			continue
		}

		low, err = strconv.ParseFloat(row[dsc.LowIndex], 64)
		if err != nil {
			continue
		}

		volume, err = strconv.ParseFloat(row[dsc.VolumeIndex], 64)
		if err != nil {
			continue
		}

		cdls = append(cdls, Candle{
			Timestamp: ts,
			High: high,
			Low: low,
			Open: open,
			Close: close,
			Volume: volume,
		})

	}
	return cdls, nil
}