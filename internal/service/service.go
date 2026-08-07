package service

import (
	"bintracker/internal/bintracker"
)

func GetBins(cfg bintracker.Config) []bintracker.Bin {
	return getBins(cfg)
}
