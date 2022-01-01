package main

var (
	// Arguments stores the arguments passed to the optimized functions.
	Arguments = map[InstructionDepth](*Stack64){}
	// cachedCalls stores the function results.
	cachedCalls = map[InstructionDepth](map[ValueWidth]ValueWidth){}
	// LastFunctionCalled is the last function called.
	LastFunctionCalled InstructionDepth
	// optimizedFunctions looks optimized instructions.
	optimizedFunctions = map[InstructionDepth]interface{}{}
)

// optimizerInitialize initialized the arguments stack.
func optimizerInitialize() {
	for f := range optimizedFunctions {
		Arguments[f] = &Stack64{}
		Arguments[f].Init(STACK_DEPTH)
	}
}

// cacheFunctionCall sees if we have a result stored.
func cacheFunctionCall(labelNum InstructionDepth) bool {
	// Skip if the function is not marked for optimizations.
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

// optimizerAnalyzeCallStack analyzes the call stack to store arguments.
func optimizerAnalyzeCallStack() {
	// Skip if the function is not marked for optimizations.
	if _, ok := optimizedFunctions[LastFunctionCalled]; !ok {
		return
	}
	Arguments[LastFunctionCalled].Push(Stack.Peek())
	if _, ok := cachedCalls[LastFunctionCalled]; !ok {
		cachedCalls[LastFunctionCalled] = map[ValueWidth]ValueWidth{}
	}
}

// optimizerAnalyzeReturnStack maps arguments to returns of pure functions.
func optimizerAnalyzeReturnStack() {
	// Skip if the function is not marked for optimizations.
	if _, ok := optimizedFunctions[LastFunctionCalled]; !ok {
		return
	}
	lastReturn := Stack.Peek()
	lastArgument := Arguments[LastFunctionCalled].Pop()
	cachedCalls[LastFunctionCalled][lastArgument] = lastReturn
}

// min is a helper min function.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max is a helper max function.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
