package gojen

// Sequence is an interface that holds the sequence of the template.
type Sequence interface {
	S(n string, is ...int) *sequence
	M(n string, is ...int) *sequence
	ForwardCtx(filteredNames ...string) *sequence
	When(i int, next func(sub Sequence)) *sequence
}
