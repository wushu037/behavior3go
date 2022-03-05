package core

import (
	"fmt"
	_ "fmt"

	b3 "behavior3go"
	_ "behavior3go/config"
)

type IBaseWorker interface {

	// Enter method, override this to use. It is called every time a node is
	// asked to execute, before the tick itself.
	OnEnter(tick *Tick)

	// Open method, override this to use. It is called only before the tick
	// callback and only if the not isn't closed.
	//
	// Note: a node will be closed if it returned `b3.RUNNING` in the tick.
	OnOpen(tick *Tick)

	//  Tick method, override this to use. This method must contain the real
	//  execution of node (perform a task, call children, etc.). It is called
	//  every time a node is asked to execute.
	OnTick(tick *Tick) b3.Status

	//  Close method, override this to use. This method is called after the tick
	//  callback, and only if the tick return a state different from
	//  `b3.RUNNING`.
	OnClose(tick *Tick)

	//  Exit method, override this to use. Called every time in the end of the
	//  execution.
	OnExit(tick *Tick)
}
type BaseWorker struct {
}

/**
 * Enter method, override this to use. It is called every time a node is
 * asked to execute, before the tick itself.
 */
func (this *BaseWorker) OnEnter(tick *Tick) {

}

/**
 * Open method, override this to use. It is called only before the tick
 * callback and only if the not isn't closed.
 *
 * Note: a node will be closed if it returned `b3.RUNNING` in the tick.
 */
func (this *BaseWorker) OnOpen(tick *Tick) {

}

/**
 * Tick method, override this to use. This method must contain the real
 * execution of node (perform a task, call children, etc.). It is called
 * every time a node is asked to execute.
 */
func (this *BaseWorker) OnTick(tick *Tick) b3.Status {
	fmt.Println("tick BaseWorker")
	return b3.ERROR
}

/**
 * Close method, override this to use. This method is called after the tick
 * callback, and only if the tick return a state different from
 * `b3.RUNNING`.
 */
func (this *BaseWorker) OnClose(tick *Tick) {

}

/**
 * Exit method, override this to use. Called every time in the end of the
 * execution.
 */
func (this *BaseWorker) OnExit(tick *Tick) {

}
