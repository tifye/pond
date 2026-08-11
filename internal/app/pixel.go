package app

func PixelBufferFromNoise(noise []float64, width, height int, scale float64) []byte {

	buff := make([]byte, 4*width*height)
	idx := 0

	for y := range height {
		for x := range width {
			noiseValue := noise[y*width+x]
			pixelValue := byte(noiseValue * 255)
			buff[idx+0] = pixelValue
			buff[idx+1] = pixelValue
			buff[idx+2] = pixelValue
			buff[idx+3] = 128
			idx += 4
		}
	}

	return buff
}
