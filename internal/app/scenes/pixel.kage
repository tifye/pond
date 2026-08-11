//go:build ignore

//kage:unit pixels

package main

var Resolution vec2

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	pixelSize := 3.0
	origin := imageSrc0Origin()
	local := srcPos - origin
	snapped := floor(local/pixelSize) * pixelSize // block corner
	snapped += pixelSize * 0.5                    // optional: block center, usually looks better
	return imageSrc0At(origin + snapped)
}
