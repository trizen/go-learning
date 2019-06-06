package main

import (
    "fmt"
    "sort"
)

type ByLength []string

func (s ByLength) Len() int           { return len(s) }
func (s ByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByLength) Less(i, j int) bool { return s[i] < s[j] }

type Tree map[rune]Tree // a recursive map

func traverse(tree Tree, acc string) string {
    if len(tree) == 1 {
        var key rune
        for k := range tree {
            key = k
            break
        }
        return traverse(tree[key], acc+string(key))
    }

    return acc
}

func lcp(strs ...string) string {

    tree := make(Tree)
    sort.Sort(ByLength(strs))

    for _, str := range strs {
        ref := tree
        if str == "" {
            return ""
        }
        for _, char := range str {
            if v, ok := ref[char]; ok {
                ref = v
                if len(ref) == 0 {
                    break
                }
            } else {
                r := make(Tree)
                ref[char] = r
                ref = r
            }
        }
    }

    return traverse(tree, "")
}

func is(a, b string) {
    if a != b {
        panic(a + " != " + b)
    }
    fmt.Println(a)
}

func main() {
    is(lcp("interspecies", "interstellar", "interstate"), "inters")
    is(lcp("throne", "throne"), "throne")
    is(lcp("throne", "dungeon"), "")
    is(lcp("throne", "", "throne"), "")
    is(lcp("cheese"), "cheese")
    is(lcp(""), "")
    is(lcp(), "")
    is(lcp("prefix", "suffix"), "")
    is(lcp("foo", "foobar"), "foo")
}
