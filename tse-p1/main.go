package main

import (
	"os"
	"tse-p1/candles"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const EthereumCSVPath = "../market-data/ETH_1min.csv"
var EthereumCSVDescriptor = candles.CandlesCsvDescriptor {
	TimestampIndex: 1,
	OpenIndex:      3,
	HighIndex:      4,
	LowIndex:       5,
	CloseIndex:     6,
	VolumeIndex:    7,
	ColumnCount:    8,
	Order:          candles.OldestLast,
}

func main() {
	candles, err := candles.LoadCandlesFromCsv(EthereumCSVPath, EthereumCSVDescriptor)
	if err != nil {
		panic(err)
	}

	var items []opts.LineData
	x := 0
	for i := 0; i < 1800000; i++ {
		if i % 1000 == 0 {
			items = append(items, opts.LineData{Value: []interface{}{x, candles[i].Close}})
			x += 1
		}
	}

	linegraph := charts.NewLine()
	linegraph.AddSeries("Close Price", items)

	f, err := os.Create("linegraph.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = linegraph.Render(f)
	if err != nil {
		panic(err)
	}
}