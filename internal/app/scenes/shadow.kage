//go:build ignore

//kage:unit pixels

package main

var Offset vec2

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
    local := srcPos - imageSrc0Origin()
    local += Offset

    c := imageSrc0At(local)

    return vec4(0.0, 0.0, 0.0, c.w*0.8)
}