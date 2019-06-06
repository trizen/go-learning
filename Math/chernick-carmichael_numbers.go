package main

// https://rosettacode.org/wiki/Chernick%27s_Carmichael_numbers

import (
    "fmt"
    "math/big"
)

var (
    fact = new(big.Int)
)

func primalityPretest(k uint64) bool {

    if (k% 3) == 0 || (k% 5) == 0 || (k% 7) == 0 || (k%11) == 0 ||
       (k%13) == 0 || (k%17) == 0 || (k%19) == 0 || (k%23) == 0 {
        return (k <= 23);
    }

    return true
}

func isChernickCarmichael(n, m uint64) bool {

    if !primalityPretest(6*m + 1) {
        return false
    }

    if !primalityPretest(12*m + 1) {
        return false
    }

    for i := uint64(1); i <= n-2; i++ {
        if !primalityPretest(((9 * m) << i) + 1) {
            return false
        }
    }

    fact.SetUint64(6*m + 1)
    if !fact.ProbablyPrime(0) {
        return false
    }

    fact.SetUint64(12*m + 1)
    if !fact.ProbablyPrime(0) {
        return false
    }

    for i := uint64(1); i <= n-2; i++ {
        fact.SetUint64(((9 * m) << i) + 1)
        if !fact.ProbablyPrime(0) {
            return false
        }
    }
    return true
}

func ccNumbers(start, end uint64) {
    for n := start; n <= end; n++ {

        multiplier := uint64(1)

        if n > 4 {
            multiplier = 1 << (n - 4)
        }

        if n > 5 {
            multiplier *= 5
        }

        for k := uint64(1); ; k++ {

            m := k * multiplier

            if isChernickCarmichael(n, m) {
                fmt.Printf("a(%d) with m = %d\n", n, m)
                break
            }
        }
    }
}

func main() {
    ccNumbers(3, 10)
}
