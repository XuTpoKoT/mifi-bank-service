package service

import "math"

func CalculateAnnuity(
	principal float64,
	annualRate float64,
	months int,
) float64 {

	r := annualRate / 100 / 12

	payment :=
		principal *
			(r * math.Pow(1+r, float64(months))) /
			(math.Pow(1+r, float64(months)) - 1)

	return math.Round(payment*100) / 100
}
