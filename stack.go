package main

type Stack32 struct {
	Array   []InstructionDepth
	Pointer int
}

func (s *Stack32) Init(depth int) {
	s.Array = make([]InstructionDepth, 0, depth)
}

func (s *Stack32) Push(val InstructionDepth) {
	if len(s.Array) == s.Pointer {
		s.Array = append(s.Array, val)
	} else {
		s.Array[s.Pointer] = val
	}
	s.Pointer++
}

func (s *Stack32) Pop() InstructionDepth {
	s.Pointer--
	return s.Array[s.Pointer]
}

type Stack64 struct {
	Array   []ValueWidth
	Pointer int
}

func (s *Stack64) Init(depth int) {
	s.Array = make([]ValueWidth, 0, depth)
}

func (s *Stack64) Push(val ValueWidth) {
	if len(s.Array) == s.Pointer {
		s.Array = append(s.Array, val)
	} else {
		s.Array[s.Pointer] = val
	}
	s.Pointer++
}

func (s *Stack64) Pop() ValueWidth {
	s.Pointer--
	return s.Array[s.Pointer]
}
