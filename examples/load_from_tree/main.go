/*
从导出的树文件加载
*/
package main

import (
	b3 "behavior3go"
	. "behavior3go/config"
	. "behavior3go/core"
	. "behavior3go/examples/share"
	. "behavior3go/loader"
	"fmt"
)

func main() {
	treeConfig, ok := LoadTreeCfg("examples/load_from_tree/tree.json")
	if !ok {
		fmt.Println("LoadTreeCfg err")
		return
	}
	//自定义节点注册
	maps := b3.NewRegisterStructMaps()
	maps.Register("Log", new(LogTest))

	//载入
	tree := CreateBevTreeFromConfig(treeConfig, maps)
	tree.Print()

	//输入板
	board := NewBlackboard()
	//循环每一帧
	for i := 0; i < 5; i++ {
		tree.Tick(i, board)
	}

}
