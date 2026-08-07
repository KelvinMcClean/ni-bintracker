package service

import (
	"bintracker/internal/bintracker"
	"bintracker/internal/gateway"
)

func getBins(cfg bintracker.Config) []bintracker.Bin {
	payload := gateway.GetBins(cfg)
	bins := []bintracker.Bin{}
	bins = appendBinsFromResponse(bins, payload.LastWeek)
	bins = appendBinsFromResponse(bins, payload.ThisWeek)
	bins = appendBinsFromResponse(bins, payload.NextWeek)
	return bins
}

func appendBinsFromResponse(acc []bintracker.Bin, week []gateway.BinDay) []bintracker.Bin {
	for _, day := range week {
		for _, bin := range day.Bins {
			acc = append(acc, bintracker.Bin{
				Date:          day.Date,
				DayOfWeekName: day.DayOfWeekName,
				ColourLabel:   bin.ColourLabel,
				Colour:        bin.Colour,
				Name:          bin.Name,
				Capacity:      bin.Capacity,
				BinType:       bin.BinType,
				UPRN:          bin.UPRN,
			})
		}
	}
	return acc
}
