package main

import (
	"./GoMerge"
	"bufio"
	"os"
	"fmt"
)
var manager *GoMerge.Manager
func main(){
	manager=new(GoMerge.Manager)
	manager.Init()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("输入指令:")
	for scanner.Scan() {
		line := scanner.Text()
		manager.CmdRout(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
