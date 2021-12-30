package main

const (
	STACK_ANALYSIS_DEPTH = 5
)

var (
	lastStackSnapshot   = make([]ValueWidth, STACK_ANALYSIS_DEPTH)
	returnStackSnapshot = make([]ValueWidth, STACK_ANALYSIS_DEPTH)
	lastStackPointer    int
	returnStackPointer  int

	cachedCalls = map[InstructionDepth]ValueWidth{}
)

func getStackSnapshot() {
	lastStackPointer = Stack.Pointer
	copy(lastStackSnapshot, Stack.Array[max(0, lastStackPointer-STACK_ANALYSIS_DEPTH):lastStackPointer])
}

func processReturnStack() {
	returnStackPointer = Stack.Pointer
	copy(returnStackSnapshot, Stack.Array[max(0, returnStackPointer-STACK_ANALYSIS_DEPTH):returnStackPointer])
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
