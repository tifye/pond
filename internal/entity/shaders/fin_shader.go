//go:build ignore

//kage:unit pixels

package main


func circle(cur vec2, target vec2, r float) vec3 {
    b := length(cur-target)
    a := step(r, b)
    return vec3(a,a,a)
}

var Time float


func Fragment(dstPos vec4, srcPos vec2, nn vec4) vec4 {
    uv := (dstPos.xy - imageDstOrigin()) / imageDstSize()

    beta := sin(uv.x*3.5-0.25)+0.5
    beta = mix(0.0, 0.2, beta)
    beta = (1.0 - step(0.7, uv.y)) * beta
    uv.y -= beta

    a :=  1.0 - circle(vec2(uv.x*0.5, uv.y), vec2(0.0,0.5), 0.5)
    return vec4(a, a.x)

    return vec4(beta, beta, beta, 1.0)
}