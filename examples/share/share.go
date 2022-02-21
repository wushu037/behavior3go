package share

import (
	b3 "behavior3go"
	"fmt"
	//. "behavior3go/actions"
	//. "behavior3go/composites"
	. "behavior3go/config"
	. "behavior3go/core"
	//. "behavior3go/decorators"
)

//自定义action节点
type LogTest struct {
	Action
	info string
}

func (this *LogTest) Initialize(setting *BTNodeCfg) {
	this.Action.Initialize(setting)
	this.info = setting.GetPropertyAsString("info")
}

func (this *LogTest) OnTick(tick *Tick) b3.Status {
	fmt.Println("logtest:",tick.GetLastSubTree(), this.info)
	return b3.SUCCESS
}
