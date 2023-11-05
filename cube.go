package cube

import (
	"fmt"
	"math"
	"time"
)

var (
	A, B, C float64

	cubeWidth float64 = 20
	width     int     = 160
	height    int     = 44

	zBuffer = make([]float64, 160*44)
	buffer  = make([]byte, 160*44)

	backgroundASCIICode int = '.'
	distanceFromCam     int = 100
	horizontalOffset    float64
	K1                  float64 = 40

	incrementSpeed float64 = 0.6
)

func calculateX(i, j, k int) float64 {
	return float64(j)*math.Sin(A)*math.Sin(B)*math.Cos(C) - float64(k)*math.Cos(A)*math.Sin(B)*math.Cos(C) +
		float64(j)*math.Cos(A)*math.Sin(C) + float64(k)*math.Sin(A)*math.Sin(C) + float64(i)*math.Cos(B)*math.Cos(C)
}

func calculateY(i, j, k int) float64 {
	return float64(j)*math.Cos(A)*math.Cos(C) + float64(k)*math.Sin(A)*math.Cos(C) -
		float64(j)*math.Sin(A)*math.Sin(B)*math.Sin(C) + float64(k)*math.Cos(A)*math.Sin(B)*math.Sin(C) -
		float64(i)*math.Cos(B)*math.Sin(C)
}

func calculateZ(i, j, k int) float64 {
	return float64(k)*math.Cos(A)*math.Cos(B) - float64(j)*math.Sin(A)*math.Cos(B) + float64(i)*math.Sin(B)
}

func calculateForSurface(cubeX, cubeY, cubeZ float32, ch int) {
	x := calculateX(int(cubeX), int(cubeY), int(cubeZ))
	y := calculateY(int(cubeX), int(cubeY), int(cubeZ))
	z := calculateZ(int(cubeX), int(cubeY), int(cubeZ)) + float64(distanceFromCam)

	ooz := 1 / z

	xp := int(float64(width)/2 + horizontalOffset + K1*ooz*x*2)
	yp := int(float64(height)/2 + K1*ooz*y)

	idx := xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = byte(ch)
		}
	}
}

// Cube 3D ASCII rotation
func Cube() {
	fmt.Print("\033[2J")
	for {
		// Clear screen
		fmt.Print("\033[H")

		buffer = make([]byte, width*height)
		zBuffer = make([]float64, width*height)
		for i := range buffer {
			buffer[i] = byte(backgroundASCIICode)
		}

		horizontalOffset = float64(-2 * cubeWidth)

		for cubeX := float32(-cubeWidth); cubeX < float32(cubeWidth); cubeX += float32(incrementSpeed) {
			for cubeY := float32(-cubeWidth); cubeY < float32(cubeWidth); cubeY += float32(incrementSpeed) {
				calculateForSurface(cubeX, cubeY, float32(-cubeWidth), '@')
				calculateForSurface(float32(cubeWidth), cubeY, cubeX, '$')
				calculateForSurface(float32(-cubeWidth), cubeY, -cubeX, '~')
				calculateForSurface(-cubeX, cubeY, float32(cubeWidth), '#')
				calculateForSurface(cubeX, float32(-cubeWidth), -cubeY, ';')
				calculateForSurface(cubeX, float32(cubeWidth), cubeY, '+')
			}
		}

		fmt.Print("\x1b[H")
		for k := 0; k < width*height; k++ {
			if k%width != 0 {
				fmt.Printf("%c", buffer[k])
			} else {
				fmt.Printf("%c", 10)
			}
		}

		A += 0.05
		B += 0.05
		C += 0.01
		time.Sleep(8 * time.Millisecond * 2)
	}
}
