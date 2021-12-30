package main

const (
	STACK_ANALYSIS_DEPTH = 5
)

var (
	Arguments          = &Stack64{}
	cachedCalls        = map[InstructionDepth](map[ValueWidth]ValueWidth){}
	LastFunctionCalled InstructionDepth
)

func init() {
	Arguments.Init(10000)
}

func cacheFunctionCall(labelNum InstructionDepth) bool {
	if _, ok := optimizedFunctions[labelNum]; !ok {
		return false
	}
	f, functionCached := cachedCalls[labelNum]
	if functionCached {
		arg, argumentCached := f[Stack.Peek()]
		if argumentCached {
			Stack.Pop()
			Stack.Push(arg)
			return true
		}
	}
	return false
}

func optimizerAnalyzeCallStack() {
	if _, ok := optimizedFunctions[LastFunctionCalled]; !ok {
		return
	}
	Arguments.Push(Stack.Peek())
	if _, ok := cachedCalls[LastFunctionCalled]; !ok {
		cachedCalls[LastFunctionCalled] = map[ValueWidth]ValueWidth{}
	}
}

func optimizerAnalyzeReturnStack() {
	if _, ok := optimizedFunctions[LastFunctionCalled]; !ok {
		return
	}
	lastReturn := Stack.Peek()
	lastArgument := Arguments.Pop()
	cachedCalls[LastFunctionCalled][lastArgument] = lastReturn
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
