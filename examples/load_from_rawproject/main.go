/*
从原生工程文件加载
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
	projectConfig, ok := LoadRawProjectCfg("examples/load_from_rawproject/example.b3")
	if !ok {
		fmt.Println("LoadRawProjectCfg err")
		return
	}

	//自定义节点注册
	maps := b3.NewRegisterStructMaps()
	maps.Register("Log", new(LogTest))

	var firstTree *BehaviorTree
	//载入
	for _, v := range projectConfig.Data.Trees {
		tree := CreateBevTreeFromConfig(&v, maps)
		tree.Print()
		if firstTree == nil {
			firstTree = tree
		}
	}

	//输入板
	board := NewBlackboard()
	//循环每一帧
	for i := 0; i < 5; i++ {
		firstTree.Tick(i, board)
	}
}
