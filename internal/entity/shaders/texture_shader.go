//go:build ignore

//kage:unit pixels

package main

var Time float
var Cursor vec2

func Fragment(dstPos vec4, srcPos vec2, color vec4, custom vec4) vec4 {
	uv := custom.xy
	return imageSrc0At(uv*imageSrc0Size() + imageSrc0Origin())
}
