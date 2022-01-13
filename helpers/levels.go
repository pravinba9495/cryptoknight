package helpers

func CalculateResistanceLevels(points []float64, before uint64, after uint64) []float64 {
	resistances := make([]float64, 0)
	for index, point := range points {
		bool := IsResistance(points, index, before, after)
		if bool {
			resistances = append(resistances, point)
		}
	}
	return resistances
}

func CalculateSupportLevels(points []float64, before uint64, after uint64) []float64 {
	supports := make([]float64, 0)
	for index, point := range points {
		bool := IsSupport(points, index, before, after)
		if bool {
			supports = append(supports, point)
		}
	}
	return supports
}

func IsSupport(points []float64, index int, before uint64, after uint64) bool {
	var aF, bF uint64

	for i := 0; i < int(before) && index-1-i >= 0; i++ {
		if points[index-1-i] > points[index] {
			bF += 1
		}
	}
	for i := 0; i < int(after) && index+1+i < len(points); i++ {
		if points[index+1+i] > points[index] {
			aF += 1
		}
	}

	return aF >= after && bF >= before
}

func IsResistance(points []float64, index int, before uint64, after uint64) bool {
	var aF, bF uint64

	for i := 0; i < int(before) && index-1-i >= 0; i++ {
		if points[index-1-i] < points[index] {
			bF += 1
		}
	}
	for i := 0; i < int(after) && index+1+i < len(points); i++ {
		if points[index+1+i] < points[index] {
			aF += 1
		}
	}

	return aF >= after && bF >= before
}
