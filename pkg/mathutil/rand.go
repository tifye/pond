package mathutil

import "math/rand/v2"

var GlobalRand = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
