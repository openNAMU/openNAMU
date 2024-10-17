package main

import (
    "fmt"

    "github.com/3th1nk/cidr"
)

func main() {
    c, err := cidr.Parse("1234")
    if err != nil {
        fmt.Println(err)
    } else if c.Contains("test") {
        fmt.Println("test")
    }
}
