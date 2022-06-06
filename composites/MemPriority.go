package composites

import (
	b3 "behavior3go"
	. "behavior3go/core"
)

type MemPriority struct {
	Composite
}

/**
 * Open method.
 * @method open
 * @param {b3.Tick} tick A tick instance.
**/
func (this *MemPriority) OnOpen(tick *Tick) {
	// 设置默认值0，OnTick时取出从0开始执行
	tick.Blackboard.Set("runningChild", 0, tick.GetTree().GetID(), this.GetID())
}

/**
 * Tick method.
 * @method tick
 * @param {b3.Tick} tick A tick instance.
 * @return {Constant} A state constant.
**/
func (this *MemPriority) OnTick(tick *Tick) b3.Status {
	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), this.GetID())
	for i := child; i < this.GetChildCount(); i++ {
		var status = this.GetChild(i).Execute(tick)

		if status != b3.FAILURE { // i note：如果是status==ERROR，也会继续下一次循环执行下一个节点，这不对吧(或者我没有理解ERROR的用法？我以为出现ERROR就终止树运行)
			if status == b3.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.GetTree().GetID(), this.GetID())
			}

			return status
		}
	}
	return b3.FAILURE
}
