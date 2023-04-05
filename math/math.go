package math

import "math"

const (
	earthRadius float64 = 6371.01 // 地球半徑（單位：公里）
)

// Radians 將角度轉換為弧度
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// CalculateDistance 計算兩個坐標之間的距離（單位：公里）
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	deltaLat := lat2 - lat1
	deltaLng := lng2 - lng1
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
