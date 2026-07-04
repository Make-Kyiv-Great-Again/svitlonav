package valhalla

func DecodePolyline(encoded string, precision int) [][2]float64 {
	var index, lat, lng int

	factor := 1.0
	for i := 0; i < precision; i++ {
		factor *= 10
	}

	coords := make([][2]float64, 0, len(encoded)/4)

	for index < len(encoded) {
		lat += decodeValue(encoded, &index)
		lng += decodeValue(encoded, &index)
		coords = append(coords, [2]float64{
			float64(lat) / factor,
			float64(lng) / factor,
		})
	}

	return coords
}

func decodeValue(encoded string, index *int) int {
	var shift, result int
	for {
		b := int(encoded[*index]) - 63
		*index++
		result |= (b & 0x1f) << shift
		shift += 5
		if b < 0x20 {
			break
		}
	}
	if result&1 != 0 {
		return ^(result >> 1)
	}
	return result >> 1
}

func MergeSegments(segments [][][2]float64) [][2]float64 {
	if len(segments) == 0 {
		return nil
	}

	merged := append([][2]float64{}, segments[0]...)
	for _, seg := range segments[1:] {
		if len(seg) == 0 {
			continue
		}
		merged = append(merged, seg[1:]...)
	}

	return merged
}
