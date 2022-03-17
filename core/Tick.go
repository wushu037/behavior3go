package core

import (
	_ "fmt"
)

/**
 * A new Tick object is instantiated every tick by BehaviorTree. It is passed
 * as parameter to the nodes through the tree during the traversal.
 *
 * The role of the Tick class is to store the instances of tree, debug,
 * target and blackboard. So, all nodes can access these informations.
 *
 * For internal uses, the Tick also is useful to store the open node after
 * the tick signal, in order to let `BehaviorTree` to keep track and close
 * them when necessary.
 *
 * This class also makes a bridge between nodes and the debug, passing the
 * node state to the debug if the last is provided.
 *
 * @module b3
 * @class Tick
翻译
BehaviorTree会在每次tick中实例化一个新的Tick对象，它在遍历过程中作为参数通过树传递给节点。
Tick类的作用是存储tree、debug、target和blackboard的实例。因此，所有节点都可以访问这些信息。
对于内部使用，Tick 也可用于在 tick 信号之后存储打开的节点，以便让 `BehaviorTree` 跟踪并在必要时关闭它们。
此类还在节点和调试之间架起了一座桥梁，如果提供了最后一个，则将节点状态传递给调试。
**/
type Tick struct {
	// The tree reference.
	tree *BehaviorTree
	
	// The debug reference.
	debug interface{}

	// The target object reference.
	target interface{}

	// The blackboard reference
	Blackboard *Blackboard

	// The list of open nodes. Update during the tree traversal
	_openNodes []IBaseNode

	// The list of open subtree node.
	// push subtree node before execute subtree.
	// pop subtree node after execute subtree.
	_openSubtreeNodes []*SubTree

	// The number of nodes entered during the tick. Update during the tree
	// traversal.
	_nodeCount int
}

func NewTick() *Tick {
	tick := &Tick{}
	tick.Initialize()
	return tick
}

/**
 * Initialization method.
 * @method Initialize
 * @construCtor
**/
func (this *Tick) Initialize() {
	// set by BehaviorTree
	this.tree = nil
	this.debug = nil
	this.target = nil
	this.Blackboard = nil

	// updated during the tick signal
	this._openNodes = nil
	this._openSubtreeNodes = nil
	this._nodeCount = 0
}

func (this *Tick) GetTree() *BehaviorTree {
	return this.tree
}

/**
 * Called when entering a node (called by BaseNode).
 * @method _enterNode
 * @param {Object} node The node that called this method.
 * @protected
**/
func (this *Tick) _enterNode(node IBaseNode) {
	this._nodeCount++
	this._openNodes = append(this._openNodes, node)

	// TODO: call debug here
}

/**
 * Callback when opening a node (called by BaseNode).
 * @method _openNode
 * @param {Object} node The node that called this method.
 * @protected
**/
func (this *Tick) _openNode(node *BaseNode) {
	// TODO: call debug here
}

/**
 * Callback when ticking a node (called by BaseNode).
 * @method _tickNode
 * @param {Object} node The node that called this method.
 * @protected
**/
func (this *Tick) _tickNode(node *BaseNode) {
	// TODO: call debug here
	//fmt.Println("Tick _tickNode :", this.debug, " id:", node.GetID(), node.GetTitle())
}

/**
 * Callback when closing a node (called by BaseNode).
 * @method _closeNode
 * @param {Object} node The node that called this method.
 * @protected
**/
func (this *Tick) _closeNode(node *BaseNode) {
	// TODO: call debug here

	ulen := len(this._openNodes)
	if ulen > 0 {
		this._openNodes = this._openNodes[:ulen-1]
	}

}

func (this *Tick) pushSubtreeNode(node *SubTree) {
	this._openSubtreeNodes = append(this._openSubtreeNodes, node)
}
func (this *Tick) popSubtreeNode() {
	ulen := len(this._openSubtreeNodes)
	if ulen > 0 {
		this._openSubtreeNodes = this._openSubtreeNodes[:ulen-1]
	}
}

/**
 * return top subtree node.
 * return nil when it is runing at major tree
 *
**/
func (this *Tick) GetLastSubTree() *SubTree {
	ulen := len(this._openSubtreeNodes)
	if ulen > 0 {
		return this._openSubtreeNodes[ulen-1]
	}
	return nil
}

/**
 * Callback when exiting a node (called by BaseNode).
 * @method _exitNode
 * @param {Object} node The node that called this method.
 * @protected
**/
func (this *Tick) _exitNode(node *BaseNode) {
	// TODO: call debug here
}

func (this *Tick) GetTarget() interface{} {
	return this.target
}
