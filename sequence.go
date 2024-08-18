package gojen

// sequence is a struct that holds the sequence of the template.
type sequence struct {
	branch bool

	n          string
	is         []int
	forwardCtx *[]string
	next       *sequence
	root       *sequence
	when       map[int]*sequence
}

// S returns a new sequence.
func S(n string, is ...int) *sequence {
	s := &sequence{
		n:  n,
		is: is,
	}

	s.root = s

	return s
}

func (s *sequence) filter(els []*E) []*E {
	res := make([]*E, 0)
	for i := range els {
		if contains(s.is, i+1) {
			res = append(res, els[i])
		}

	}

	return res
}

// S adds new sequence to the chain.
func (s *sequence) S(n string, is ...int) *sequence {
	next := &sequence{
		n:    n,
		is:   is,
		root: s.root,
	}

	if len(s.when) > 0 {
		for _, w := range s.when {
			w.next = next
		}
	} else {
		s.next = next
	}

	return next
}

// When adds a condition to the sequence.
func (s *sequence) When(i int, thenFn func(sub Sequence)) *sequence {
	if s.when == nil {
		s.when = map[int]*sequence{}
	}
	s.when[i] = &sequence{
		branch: true,
		root:   s.root,
	}

	if thenFn != nil {
		thenFn(s.when[i])
	}

	return s
}

// ForwardCtx forward current context to the next sequence to inherit.
func (s *sequence) ForwardCtx(filteredNames ...string) *sequence {
	s.forwardCtx = &filteredNames
	return s
}
