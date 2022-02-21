package actions

import (
	"fmt"

	b3 "behavior3go"
	. "behavior3go/config"
	. "behavior3go/core"
)

type Log struct {
	Action
	info string
}

func (this *Log) Initialize(setting *BTNodeCfg) {
	this.Action.Initialize(setting)
	this.info = setting.GetPropertyAsString("info")
}

func (this *Log) OnTick(tick *Tick) b3.Status {
	fmt.Println("log:", this.info)
	return b3.SUCCESS
}
