package main

// Stack32 is a stack for uint32.
type Stack32 struct {
	// Array is the internal of stack.
	Array []InstructionDepth
	// Pointer is the Array pointer such that
	// in order to avoid constant reallocations,
	// keep the memory but move the pointer.
	Pointer int
}

// Init slab allocates the underlying array with
// some initial passed capacity.
func (s *Stack32) Init(depth int) {
	s.Array = make([]InstructionDepth, 0, depth)
}

// Push pushes a value to the stack.
func (s *Stack32) Push(val InstructionDepth) {
	if len(s.Array) == s.Pointer {
		s.Array = append(s.Array, val)
	} else {
		s.Array[s.Pointer] = val
	}
	s.Pointer++
}

// Peek returns the top of the stack.
func (s *Stack32) Peek() InstructionDepth {
	return s.Array[s.Pointer-1]
}

// Pop pops the top of the stack and returns it.
func (s *Stack32) Pop() InstructionDepth {
	s.Pointer--
	return s.Array[s.Pointer]
}

// Stack32 is a stack for uint64.
type Stack64 struct {
	// Array is the internal of stack.
	Array []ValueWidth
	// Pointer is the Array pointer such that
	// in order to avoid constant reallocations,
	// keep the memory but move the pointer.
	Pointer int
}

// Init slab allocates the underlying array with
// some initial passed capacity.
func (s *Stack64) Init(depth int) {
	s.Array = make([]ValueWidth, 0, depth)
}

// Push pushes a value to the stack.
func (s *Stack64) Push(val ValueWidth) {
	if len(s.Array) == s.Pointer {
		s.Array = append(s.Array, val)
	} else {
		s.Array[s.Pointer] = val
	}
	s.Pointer++
}

// Peek returns the top of the stack.
func (s *Stack64) Peek() ValueWidth {
	return s.Array[s.Pointer-1]
}

// Pop pops the top of the stack and returns it.
func (s *Stack64) Pop() ValueWidth {
	s.Pointer--
	return s.Array[s.Pointer]
}
