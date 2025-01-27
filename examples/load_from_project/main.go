/*
从导出的工程文件加载
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
	projectConfig, ok := LoadProjectCfg("examples/load_from_project/project.json")
	if !ok {
		fmt.Println("LoadTreeCfg err")
		return
	}

	//自定义节点注册
	maps := b3.NewRegisterStructMaps()
	maps.Register("Log", new(LogTest))

	var firstTree *BehaviorTree
	//载入
	for _, v := range projectConfig.Trees {
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
