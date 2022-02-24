package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//工程json类型
type BTProjectCfg struct {
	ID     string      `json:"id"`
	Select string      `json:"selectedTree"`
	Scope  string      `json:"scope"`
	Trees  []BTTreeCfg `json:"trees"`
}

//加载
func LoadProjectCfg(relativePath string) (*BTProjectCfg, bool) {

	wdPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := fmt.Sprintf("%s/%s", wdPath, relativePath)

	var project BTProjectCfg
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("LoadProjectCfg fail:", err)
		return nil, false
	}
	err = json.Unmarshal(file, &project)
	if err != nil {
		fmt.Println("LoadProjectCfg fail, ummarshal:", err, len(file))
		return nil, false
	}

	//fmt.Println("load tree:", tree.Title, " nodes:", len(tree.Nodes))
	return &project, true
}
