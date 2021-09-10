package main

import "math"

func calcDistancia(lati, longi, latf, longf float64) float64 {
	cos := math.Cos(rad(latf - lati))
	topCos := math.Cos(rad(latf)) * math.Cos(rad(lati)) * (1 - math.Cos(rad(longf-longi)))
	return math.Acos(cos-topCos) * raio
}

func rad(deg float64) float64 {
	return deg * math.Pi / 180
}
