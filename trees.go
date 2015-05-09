package main

import (
    "fmt"
    "strings"
)

type Node struct {
    char rune
    link *Node
    last bool
}

func lcp(strs ...string) []*Node {
    var trees []*Node

    for _, str := range strs {
        var hash *Node
        hash = &Node{0, hash, false}

        ref := hash
        for _, char := range str {
            var node *Node
            ref.char = char
            ref.link = &Node{char, node, false}
            ref = ref.link
        }

        ref.last = true
        trees = append(trees, hash)
    }

    return trees
}

func traverse(tree *Node, spaces int) {
    if !tree.last {
        fmt.Println(strings.Repeat(" ", spaces), string(tree.char))
        traverse(tree.link, spaces+4)
    }
}

func main() {
    trees := lcp("interspecies", "interstellar", "interstate")
    for _, tree := range trees {
        traverse(tree, 0)
    }
}
