package core

import (
	"fmt"

	b3 "behavior3go"
	"behavior3go/config"
)

/**
 * The BehaviorTree class, as the name implies, represents the Behavior Tree
 * structure.
 *
 * There are two ways to construct a Behavior Tree: by manually setting the
 * root node, or by loading it from a data structure (which can be loaded
 * from a JSON). Both methods are shown in the examples below and better
 * explained in the user guide.
 *
 * The tick method must be called periodically, in order to send the tick
 * signal to all nodes in the tree, starting from the root. The method
 * `BehaviorTree.tick` receives a target object and a blackboard as
 * parameters. The target object can be anything: a game agent, a system, a
 * DOM object, etc. This target is not used by any piece of Behavior3JS,
 * i.e., the target object will only be used by custom nodes.
 *
 * The blackboard is obligatory and must be an instance of `Blackboard`. This
 * requirement is necessary due to the fact that neither `BehaviorTree` or
 * any node will store the execution variables in its own object (e.g., the
 * BT does not store the target, information about opened nodes or number of
 * times the tree was called). But because of this, you only need a single
 * tree instance to control multiple (maybe hundreds) objects.
 *
 * Manual construction of a Behavior Tree
 * --------------------------------------
 *
 *     var tree = new b3.BehaviorTree();
 *
 *     tree.root = new b3.Sequence({children:[
 *       new b3.Priority({children:[
 *         new MyCustomNode(),
 *         new MyCustomNode()
 *       ]}),
 *       ...
 *     ]});
 *
 *
 * Loading a Behavior Tree from data structure
 * -------------------------------------------
 *
 *     var tree = new b3.BehaviorTree();
 *
 *     tree.load({
 *       'title'       : 'Behavior Tree title'
 *       'description' : 'My description'
 *       'root'        : 'node-id-1'
 *       'nodes'       : {
 *         'node-id-1' : {
 *           'name'        : 'Priority', // this is the node type
 *           'title'       : 'Root Node',
 *           'description' : 'Description',
 *           'children'    : ['node-id-2', 'node-id-3'],
 *         },
 *         ...
 *       }
 *     })
 *
 *
 * @module b3
 * @class BehaviorTree
**/
type BehaviorTree struct {

	// The tree id, must be unique. By default, created with `b3.createUUID`.
	id string

	// The tree title
	title string

	// Description of the tree
	description string

	// A dictionary with (key-value) properties. Useful to define custom
	// variables in the visual editor.
	properties map[string]interface{}

	// The reference to the root node. Must be an instance of `b3.BaseNode`.
	root IBaseNode

	// The reference to the debug instance
	debug interface{}

	dumpInfo *config.BTTreeCfg
}

func NewBeTree() *BehaviorTree {
	tree := &BehaviorTree{}
	tree.Initialize()
	return tree
}

/**
 * Initialization method.
 * @method Initialize
 * @construCtor
**/
func (this *BehaviorTree) Initialize() {
	this.id = b3.CreateUUID()
	this.title = "The behavior tree"
	this.description = "Default description"
	this.properties = make(map[string]interface{})
	this.root = nil
	this.debug = nil
}

func (this *BehaviorTree) GetID() string {
	return this.id
}

func (this *BehaviorTree) GetTitile() string {
	return this.title
}

func (this *BehaviorTree) SetDebug(debug interface{}) {
	this.debug = debug
}

func (this *BehaviorTree) GetRoot() IBaseNode {
	return this.root
}

/**
 * This method loads a Behavior Tree from a data structure, populating this
 * object with the provided data. Notice that, the data structure must
 * follow the format specified by Behavior3JS. Consult the guide to know
 * more about this format.
 *
 * You probably want to use custom nodes in your BTs, thus, you need to
 * provide the `names` object, in which this method can find the nodes by
 * `names[NODE_NAME]`. This variable can be a namespace or a dictionary,
 * as long as this method can find the node by its name, for example:
 *
 *     //json
 *     ...
 *     'node1': {
 *       'name': MyCustomNode,
 *       'title': ...
 *     }
 *     ...
 *
 *     //code
 *     var bt = new b3.BehaviorTree();
 *     bt.load(data, {'MyCustomNode':MyCustomNode})
 *
 *
 * @method load
 * @param {Object} data The data structure representing a Behavior Tree.
 * @param {Object} [names] A namespace or dict containing custom nodes.
**/
func (this *BehaviorTree) Load(data *config.BTTreeCfg, maps *b3.RegisterStructMaps, extMaps *b3.RegisterStructMaps) {
	this.title = data.Title             //|| this.title;
	this.description = data.Description // || this.description;
	this.properties = data.Properties   // || this.properties;
	this.dumpInfo = data
	nodes := make(map[string]IBaseNode)

	// Create the node list (without connection between them)

	for id, nodeCfg := range data.Nodes {
		var node IBaseNode

		if nodeCfg.Category == "tree" {
			node = new(SubTree)
		} else {
			if extMaps != nil && extMaps.CheckElem(nodeCfg.Name) {
				// Look for the name in custom nodes
				if tnode, err := extMaps.New(nodeCfg.Name); err == nil {
					node = tnode.(IBaseNode)
				}
			} else {
				if tnode, err2 := maps.New(nodeCfg.Name); err2 == nil {
					node = tnode.(IBaseNode)
				} else {
					//fmt.Println("new ", nodeCfg.Name, " err:", err2)
				}
			}
		}

		if node == nil {
			// Invalid node name
			panic("BehaviorTree.load: Invalid node name:" + nodeCfg.Name + ",title:" + nodeCfg.Title)

		}

		node.Ctor()
		node.Initialize(&nodeCfg)
		node.SetBaseNodeWorker(node.(IBaseWorker))
		nodes[id] = node
	}

	// Connect the nodes
	for id, nodeCfg := range data.Nodes {
		node := nodes[id]

		if node.GetCategory() == b3.COMPOSITE && nodeCfg.Children != nil {
			for i := 0; i < len(nodeCfg.Children); i++ {
				var cid = nodeCfg.Children[i]
				comp := node.(IComposite)
				comp.AddChild(nodes[cid])
			}
		} else if node.GetCategory() == b3.DECORATOR && len(nodeCfg.Child) > 0 {
			dec := node.(IDecorator)
			dec.SetChild(nodes[nodeCfg.Child])
		}
	}

	this.root = nodes[data.Root]
}

/**
 * This method dump the current BT into a data structure.
 *
 * Note: This method does not record the current node parameters. Thus,
 * it may not be compatible with load for now.
 *
 * @method dump
 * @return {Object} A data object representing this tree.
**/
func (this *BehaviorTree) dump() *config.BTTreeCfg {
	return this.dumpInfo
}

func (this *BehaviorTree) Print() {
	printNode(this.root, 0)
}

/**
 * Propagates the tick signal through the tree, starting from the root.
 *
 * This method receives a target object of any type (Object, Array,
 * DOMElement, whatever) and a `Blackboard` instance. The target object has
 * no use at all for all Behavior3JS components, but surely is important
 * for custom nodes. The blackboard instance is used by the tree and nodes
 * to store execution variables (e.g., last node running) and is obligatory
 * to be a `Blackboard` instance (or an object with the same interface).
 *
 * Internally, this method creates a Tick object, which will store the
 * target and the blackboard objects.
 *
 * Note: BehaviorTree stores a list of open nodes from last tick, if these
 * nodes weren't called after the current tick, this method will close them
 * automatically.
 *
 * @method tick
 * @param {Object} target A target object.
 * @param {Blackboard} blackboard An instance of blackboard object.
 * @return {Constant} The tick signal state.
 */
/**
翻译：
从root开始在tree中传播tick信号。
此方法接收任何类型的target对象（对象、数组、DOMElement等）和一个“黑板”实例。
==target对象对于所有 Behavior3 组件根本没有用处，但对于自定义节点来说肯定很重要== (wushu:或许是一个类似于"上下文"、"拓展参数"的一个人参数)
blackboard被tree和node用来存储执行变量（例如，最后一个运行的节点），并且必须是“黑板”实例（或具有相同接口的对象）

在内部，此方法创建一个 Tick 对象，该对象将存储target和blackboard对象。
注意： BehaviorTree 存储了一个从最后一个tick开始的节点列表，如果这些节点在当前tick之后没有被调用，这个方法会自动关闭它们。

@target：一个目标对象。
@blackboard：一个黑板实例
@return：滴答信号状态。
**/
func (this *BehaviorTree) Tick(target interface{}, blackboard *Blackboard) b3.Status {
	if blackboard == nil {
		panic("The blackboard parameter is obligatory and must be an instance of b3.Blackboard")
	}

	// 创建tick对象
	var tick = NewTick()
	tick.debug = this.debug
	tick.target = target
	tick.Blackboard = blackboard
	tick.tree = this

	// 执行节点逻辑。内部会按照结构顺序，调用所有节点的execute
	// 如果有running的节点
	var state = this.root._execute(tick)

	// 关闭上一次tick的节点(如果需要)
	// openNodes: 其实就是tick后处于running状态的节点；注意：一个节点处于running时，其父节点可能也处于running状态，或许会有一条"running"链
	var lastOpenNodes = blackboard._getTreeData(this.id).OpenNodes // 上一次tick的openNodes
	var currOpenNodes []IBaseNode
	currOpenNodes = append(currOpenNodes, tick._openNodes...) // 本次tick的openNodes

	// 如果在本次tick内仍处于open状态，则不会关闭
	var start = 0 // 从第几个节点开始关闭
	for i := 0; i < b3.MinInt(len(lastOpenNodes), len(currOpenNodes)); i++ {
		start = i + 1
		// 遍历本次和上次的running调用链，若存在状态不同的节点，则从这个节点之后的所有节点都要被关闭
		if lastOpenNodes[i] != currOpenNodes[i] {
			break
		}
	}

	// 关闭这个节点及其所有后续节点
	for i := len(lastOpenNodes) - 1; i >= start; i-- {
		node := lastOpenNodes[i]
		// 打印节点及该节点是否已被关闭。会发现打印结果是有规律的，每次运行的打印结果都是固定的
		fmt.Println(node)
		fmt.Println("is-open:", tick.Blackboard.Get("isOpen", tick.tree.id, node.GetID()))
		lastOpenNodes[i]._close(tick)
	}

	// todo 可运行`memsubtree/main.go`触发以下逻辑进行分析
	// 通过上面的打印分析得出：
	//  类似这样的一个树结构：一个子树st被主树的两个分支a、b调用。
	// 	若本次tick通过分支a进入st，上次tick通过分支b进入st
	//  则上次tick造成的st中的running节点就要被关闭(running节点存在于openNodes中)。
	//  这相当于为本次tick要用到st做了初始化
	//  但这种初始化的方式或许欠妥(不适用于所有场景)：
	// 		- 在主观上，引用子树就是相当于把子树的节点添加到主树中。子树起到的是对子树节点结构的封装作用
	//		- 在本引擎中，分支a、b都引用了同一个子树st，但st内节点的状态在不同分支下却都是同一个(黑板通过nodeId存放数据，不同分支引用的st中的nodeID是一样的)
	//		- 在常规需求中我们应该更想这样：a分支下的st节点和b分支下的st节点，应该是两套节点，他们的内存应该是分开的。可以给子树节点id依据分支加上不同的前缀
	//
	// 冗余的触发情况：本次tick没有openNodes，上次tick有，但上次tick的openNodes在本次tick运行时已经被正常close了，也会再触发这里的close，这显然是多余的


	// 填充黑板数据
	// 本次tick的openNodes保存到黑板中在下次tick时使用
	blackboard._getTreeData(this.id).OpenNodes = currOpenNodes
	// nodeCount：本次tick中，执行了_enter()的所有节点数量。没看到有什么用途
	blackboard.SetTree("nodeCount", tick._nodeCount, this.id)

	return state
}

func printNode(root IBaseNode, blk int) {

	//fmt.Println("new node:", root.Name, " children:", len(root.Children), " child:", root.Child)
	for i := 0; i < blk; i++ {
		fmt.Print(" ") //缩进
	}

	//fmt.Println("|—<", root.Name, ">") //打印"|—<id>"形式
	fmt.Print("|—", root.GetTitle())

	if root.GetCategory() == b3.DECORATOR {
		dec := root.(IDecorator)
		if dec.GetChild() != nil {
			//fmt.Print("=>")
			printNode(dec.GetChild(), blk+3)
		}
	}

	fmt.Println("")
	if root.GetCategory() == b3.COMPOSITE {
		comp := root.(IComposite)
		if comp.GetChildCount() > 0 {
			for i := 0; i < comp.GetChildCount(); i++ {
				printNode(comp.GetChild(i), blk+3)
			}
		}
	}

}
