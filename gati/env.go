package gati

import "os"

// RNGServAddr - RNG_SERVICE_ADDRESS, http://rng.local:2007
var RNGServAddr string

// BaseRNGURL - RNG_SERVICE_ADDRESS + /numbers
var BaseRNGURL string

func init() {
	RNGServAddr = os.Getenv("RNG_SERVICE_ADDRESS")
	BaseRNGURL = RNGServAddr + "/numbers"
}
