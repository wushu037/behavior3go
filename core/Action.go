package core

import (
	b3 "behavior3go"
	. "behavior3go/config"
)

type IAction interface {
	IBaseNode
}

/**
 * Action is the base class for all action nodes. Thus, if you want to create
 * new custom action nodes, you need to inherit from this class. For example,
 * take a look at the Runner action:
 *
 *     var Runner = b3.Class(b3.Action, {
 *       name: 'Runner',
 *
 *       tick: function(tick) {
 *         return b3.RUNNING;
 *       }
 *     });
 *
 * @module b3
 * @class Action
 * @extends BaseNode
**/
type Action struct {
	BaseNode
	BaseWorker // BaseNode 里已经有 IBaseWorker 了，为什么这里还要声明？因为IBaseWorker的值其实就是这个BaseWorker，IBaseWorker的方法最终就是这个BaseWorker提供的
}

func (this *Action) Ctor() {
	this.category = b3.ACTION
}
func (this *Action) Initialize(params *BTNodeCfg) {

	//this.id = b3.CreateUUID()
	this.BaseNode.Initialize(params)
	//this.BaseNode.IBaseWorker = this
	this.parameters = make(map[string]interface{})
	this.properties = make(map[string]interface{})
}
