//go:build ignore

//kage:unit pixels

package main

var Offset vec2

func circle(cur vec2, target vec2, r float) vec4 {
    b := length(cur-target)
    a := 1.0 - step(r, b)
    return vec4(a,a,a,a)
}


func halfCircle(cur vec2, target vec2, r float) vec4 {
    c := circle(cur, target, r)
    alpha := step(0.5, cur.x)
    return c * alpha
}

func rotate(cur vec2, v vec2, a float) vec2 {
    v = cur - vec2(0.5, 0.5)

    rotX := v.x * cos(a) - v.y * sin(a) 
    rotY := v.x * sin(a) + v.y * cos(a) 

    return vec2(rotX, rotY) + vec2(0.5, 0.5)
}

// rotate vertex positions to rotate lilypad
func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    uv := texCoord 

    uv = rotate(uv, uv, 1.2)

    angle := 3.14 / 4

    h1a := -(angle * 0.5)
    rotated := rotate(uv, uv, h1a)
    h1 := halfCircle(rotated, vec2(0.5, 0.5), 0.5)

    h2a := (angle * 0.5) + 3.14
    rotated = rotate(uv, uv, h2a)
    h2 := halfCircle(rotated, vec2(0.5, 0.5), 0.5)


    return max(h1, h2) * vec4(0.02, 0.6, 0.411, 1.0) 
}