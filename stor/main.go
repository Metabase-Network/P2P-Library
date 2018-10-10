package main

import (
	"fmt"

	"github.com/Metabase-Network/vasuki/stor"
)

func main() {
	n, e := stor.start("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	fmt.Println(n, e) 
}
