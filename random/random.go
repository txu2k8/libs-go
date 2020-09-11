package random

// Package provides various convenience functions for dealing with random numbers and strings.

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	mathrand "math/rand"
	"time"
)

// Set of characters to use for generating random strings
const (
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Numerals     = "1234567890"
	Alphanumeric = Alphabet + Numerals
	ASSIC        = Alphanumeric + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

// ErrMinMax .
var ErrMinMax = errors.New("Min cannot be greater than max")

// IntRange returns a random integer in the range from min to max.
func IntRange(min, max int) (int, error) {
	var result int
	switch {
	case min > max:
		// Fail with error
		return result, ErrMinMax
	case max == min:
		result = max
	case max > min:
		maxRand := max - min
		b, err := rand.Int(rand.Reader, big.NewInt(int64(maxRand)))
		if err != nil {
			return result, err
		}
		result = min + int(b.Int64())
	}
	return result, nil
}

// String returns a random string n characters long, composed of entities
// from charset.
func String(n int, charset string) (string, error) {
	randstr := make([]byte, n) // Random string to return
	charlen := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		b, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			return "", err
		}
		r := int(b.Int64())
		randstr[i] = charset[r]
	}
	return string(randstr), nil
}

// StringRange returns a random string at least min and no more than max
// characters long, composed of entitites from charset.
func StringRange(min, max int, charset string) (string, error) {
	//
	// First determine the length of string to be generated
	//
	var err error      // Holds errors
	var strlen int     // Length of random string to generate
	var randstr string // Random string to return
	strlen, err = IntRange(min, max)
	if err != nil {
		return randstr, err
	}
	randstr, err = String(strlen, charset)
	if err != nil {
		return randstr, err
	}
	return randstr, nil
}

// AlphaStringRange returns a random alphanumeric string at least min and no more
// than max characters long.
func AlphaStringRange(min, max int) (string, error) {
	return StringRange(min, max, Alphanumeric)
}

// AlphaString returns a random alphanumeric string n characters long.
func AlphaString(n int) (string, error) {
	return String(n, Alphanumeric)
}

// NumString returns a random Numerals string n characters long.
func NumString(n int) (string, error) {
	return String(n, Numerals)
}

// ChoiceString returns a random selection from an array of strings.
func ChoiceString(choices []string) (string, error) {
	var winner string
	length := len(choices)
	i, err := IntRange(0, length)
	winner = choices[i]
	return winner, err
}

// ChoiceInt returns a random selection from an array of integers.
func ChoiceInt(choices []int) (int, error) {
	var winner int
	length := len(choices)
	i, err := IntRange(0, length)
	winner = choices[i]
	return winner, err
}

// A WeightItem contains a generic item and a weight controlling the frequency with
// which it will be selected.
type WeightItem struct {
	Weight int
	Item   interface{}
}

// WeightedChoice used weighted random selection to return one of the supplied
// choices.  Weights of 0 are never selected.  All other weight values are
// relative.  E.g. if you have two choices both weighted 3, they will be
// returned equally often; and each will be returned 3 times as often as a
// choice weighted 1.
func WeightedChoice(choices []WeightItem) (WeightItem, error) {
	// Based on this algorithm:
	//     http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
	var ret WeightItem
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r, err := IntRange(0, sum)
	if err != nil {
		return ret, err
	}
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c, nil
		}
	}
	err = errors.New("Internal error - code should not reach this point")
	return ret, err
}

// Choice .
func Choice(seq []interface{}) interface{} {
	if len(seq) == 0 {
		return seq
	}
	idx, _ := IntRange(0, len(seq)-1)
	return seq[idx]
}

// Shuffle .
func Shuffle(x []interface{}) {
	r := mathrand.New(mathrand.NewSource(time.Now().Unix()))
	for len(x) > 0 {
		n := len(x)
		randIndex := r.Intn(n)
		x[n-1], x[randIndex] = x[randIndex], x[n-1]
		x = x[:n-1]
	}
}

// Sample Chooses k unique random elements from a population sequence or set
func Sample(population []interface{}, k int) []interface{} {
	Shuffle(population)
	min, _ := IntRange(0, len(population)-k)
	return population[min : min+k]
}

// Choices Return a k sized list of population elements chosen with replacement
// TODO
func Choices(population []interface{}, weights float64, k int) []interface{} {
	Shuffle(population)
	min, _ := IntRange(0, len(population)-k)
	return population[min : min+k]
}

// -------------------- use case  -------------------

// ChoiceStrArr Choice for []string
func ChoiceStrArr(arr []string) string {
	s := make([]interface{}, len(arr))
	for i, v := range arr {
		s[i] = v
	}
	return fmt.Sprintf("%v", Choice(s))
}

// SampleStrArr Choice for []string
func SampleStrArr(arr []string, k int) []string {
	s := make([]interface{}, len(arr))
	for i, v := range arr {
		s[i] = v
	}
	spArr := make([]string, k)
	for j, sp := range Sample(s, k) {
		spArr[j] = fmt.Sprintf("%v", sp)
	}
	return spArr
}
