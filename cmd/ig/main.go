package main

import (
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/wmw9/ig"
	"os"
)

var (
	igSessionId = os.Getenv("IG_SESSION_ID1")
)

func main() {
	fmt.Println("vim-go")
	list := make([]int, 0)
	list = append(list, 3207093470)
	//list = append(list, 2209342170)
	//list = append(list, 6186074557)
	//list = append(list, 1418686499) // tati
	//list = append(list, 1434138860) // eva
	//list = append(list, 2944757465) // mira
	//list = append(list, 192322605)  // лиззка
	stories := ig.Get(igSessionId).Stories(list)
	//stories := ig.Get(igSessionId).After(1619598794).Stories(410036941, 2209342170)
	//list := make([]string, 0)
	//	list = append(list, "olyashaasaxon")
	//list = append(list, "arina_gp")
	//posts := ig.Get(igSessionId).After(0).Posts(list)
	//pp.Println(posts)
	pp.Println(stories.Stories)
	//_ = os.WriteFile("./posts.txt", posts, 0644)

}
