package main

//~ Author: Daniel "Trizen" Șuteu
//~ License: GPLv3
//~ Date: 01 May 2015
//~ Website: http://github.com/trizen

//~ The arithmetic coding algorithm.

//~ See: http://en.wikipedia.org/wiki/Arithmetic_coding#Arithmetic_coding_as_a_generalized_change_of_radix

import (
    "fmt"
    "math/big"
    "strings"
)

func asciibet() [256]byte {
    r := [256]byte{}
    for i := 0; i < 256; i++ {
        r[i] = byte(i)
    }
    return r
}

func cumulative_freq(freq map[byte]int64) map[byte]int64 {
    total := int64(0)
    cf := make(map[byte]int64)
    for _, c := range asciibet() {
        if v, ok := freq[c]; ok {
            cf[c] = total
            total += v
        }
    }
    return cf
}

func arithmethic_coding(str string) (*big.Int, *big.Int, map[byte]int64) {

    // Convert the string into a slice of bytes
    chars := []byte(str)

    // The frequency characters
    freq := make(map[byte]int64)
    for _, c := range chars {
        freq[c] += 1
    }

    // The cumulative frequency
    cf := cumulative_freq(freq)

    // Base
    base := len(chars)

    // Lower bound
    L := big.NewInt(0)

    // Product of all frequencies
    pf := big.NewInt(1)

    // Each term is multiplied by the product of the
    // frequencies of all previously occurring symbols
    bigBase := big.NewInt(int64(base))
    bigLim := big.NewInt(int64(base - 1))

    for i := 0; i < base; i++ {

        bigI := big.NewInt(int64(i))

        diff := big.NewInt(0)
        diff.Sub(bigLim, bigI)

        pow := big.NewInt(0)
        pow.Exp(bigBase, diff, nil)

        x := big.NewInt(1)
        x.Mul(big.NewInt(cf[chars[i]]), pow)

        L.Add(L, x.Mul(x, pf))
        pf.Mul(pf, big.NewInt(freq[chars[i]]))
    }

    // Upper bound
    U := big.NewInt(0)
    U.Set(L)
    U.Add(U, pf)

    bigOne := big.NewInt(1)
    bigZero := big.NewInt(0)
    bigTen := big.NewInt(10)

    tmp := big.NewInt(0).Set(pf)
    pow10 := big.NewInt(0)

    for {
        tmp.Div(tmp, bigTen)
        if tmp.Cmp(bigZero) == 0 {
            break
        }
        pow10.Add(pow10, bigOne)
    }

    diff := big.NewInt(0)
    diff.Sub(U, bigOne)
    diff.Div(diff, big.NewInt(0).Exp(bigTen, pow10, nil))

    return diff, pow10, freq
}

func arithmethic_decoding(num *big.Int, pow *big.Int, freq map[byte]int64) string {

    pow10 := big.NewInt(10)

    enc := big.NewInt(0).Set(num)
    enc.Mul(enc, pow10.Exp(pow10, pow, nil))

    base := int64(0)
    for _, v := range freq {
        base += v
    }

    // Character range
    r := asciibet()

    // Create the cumulative frequency table
    cf := cumulative_freq(freq)

    // Create the dictionary
    dict := make(map[int64]byte)
    j := int64(0)
    for i := range r {
        c := r[i]
        if v, ok := freq[c]; ok {
            dict[j] = c
            j += v
        }
    }

    // Fill the gaps in the dictionary
    lchar := -1
    for i := int64(0); i < base; i++ {
        if v, ok := dict[i]; ok {
            lchar = int(v)
        } else if lchar != -1 {
            dict[i] = byte(lchar)
        }
    }

    // Decode the input number
    decoded := make([]byte, base)
    bigBase := big.NewInt(base)

    for i := base - 1; i >= 0; i-- {

        pow := big.NewInt(0)
        pow.Exp(bigBase, big.NewInt(i), nil)

        div := big.NewInt(0)
        div.Div(enc, pow)

        c := dict[div.Int64()]
        fv := freq[c]
        cv := cf[c]

        prod := big.NewInt(0).Mul(pow, big.NewInt(cv))
        diff := big.NewInt(0).Sub(enc, prod)
        enc.Div(diff, big.NewInt(fv))

        decoded[base-i-1] = c
    }

    // Return the decoded output
    return string(decoded)
}

func main() {
    strSlice := []string{
        `DABDDB`,
        `DABDDBBDDBA`,
        `ABBDDD`,
        `ABRACADABRA`,
        `CoMpReSSeD`,
        `Sidef`,
        `Trizen`,
        `Google`,
        `TOBEORNOTTOBEORTOBEORNOT`,
        `In a positional numeral system the radix, or base, is numerically equal to a number of different symbols ` +
            `used to express the number. For example, in the decimal system the number of symbols is 10, namely 0, 1, 2, ` +
            `3, 4, 5, 6, 7, 8, and 9. The radix is used to express any finite integer in a presumed multiplier in polynomial ` +
            `form. For example, the number 457 is actually 4×102 + 5×101 + 7×100, where base 10 is presumed but not shown explicitly.`,
    }

    for _, str := range strSlice {
        enc, pow, freq := arithmethic_coding(str)
        dec := arithmethic_decoding(enc, pow, freq)

        fmt.Println("Encoded:", enc)
        fmt.Println("Decoded:", dec)

        if str != dec {
            panic("\tHowever that is incorrect!")
        }

        fmt.Println(strings.Repeat("-", 80))
    }
}
