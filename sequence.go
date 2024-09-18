package gojen

import (
	"fmt"
	"strings"

	"github.com/cirius-go/gojen/util"
)

// Seq reresents a sequence.
type Seq struct {
	root        *Seq                     `yaml:"-"`
	DName       string                   `yaml:"d_name"`
	EName       string                   `yaml:"e_name"`
	ForwardArgs util.MapExisting[string] `yaml:"forward_args,omitempty"` // forward state for all children.
	IsCase      bool                     `yaml:"is_case,omitempty"`      // is case.
	Next        *Seq                     `yaml:"next,omitempty"`
	Cases       SeqCases                 `yaml:"cases,omitempty"`
}

// SeqCases is a map of cases.
// Key is the eName.
type SeqCases map[string]*Seq

// NewSeq creates a new sequence.
func NewSeq(dName, eName string) *Seq {
	s := &Seq{
		DName:       dName,
		EName:       eName,
		ForwardArgs: util.MapExisting[string]{},
		Cases:       SeqCases{},
	}

	s.root = s

	return s
}

// AllLast returns all the last states of the sequence.
func (s *Seq) AllLast() []*Seq {
	res := []*Seq{}

	// If the sequence has no cases, it's no braching.
	if len(s.Cases) == 0 {
		if s.Next == nil {
			res = append(res, s)

			return res
		}

		children := s.Next.AllLast()
		res = append(res, children...)

		return res
	}

	for _, c := range s.Cases {
		children := c.AllLast()
		res = append(res, children...)
	}

	return res
}

func (s *Seq) append(dName, eName string) *Seq {
	n := NewSeq(dName, eName)
	n.root = s.root
	lasts := s.AllLast()
	for _, l := range lasts {
		l.Next = n
	}
	return n
}

func (s *Seq) Append(dName, eName string, moreENames ...string) *Seq {
	n := s.append(dName, eName)
	for _, eName := range moreENames {
		n = n.append(dName, eName)
	}
	return n
}

func (s *Seq) Forward(argNames ...string) *Seq {
	for _, argName := range argNames {
		s.ForwardArgs.Add(argName)
	}
	return s
}

func (s *Seq) Select(dName string, eNames []string, handler func(ss SeqSwitcher)) *Seq {
	if len(eNames) == 0 {
		panic(fmt.Errorf("no element names provided for selection"))
	}
	for _, eName := range eNames {
		c := NewSeq(dName, eName)
		c.root = s.root
		c.IsCase = true
		if s.Cases == nil {
			s.Cases = SeqCases{}
		}
		s.Cases[eName] = c
	}

	if handler != nil {
		handler(s)
	}

	return s
}

func (s *Seq) When(eName string, handler util.PRFunc[*Seq, *Seq]) SeqSwitcher {
	if c, exists := s.Cases[eName]; exists {
		return handler(c)
	}
	return s
}

func (s *Seq) String() string {
	return strings.Join(s.root.string(""), "\n")
}

func (s *Seq) string(indent string) []string {
	res := make([]string, 0) // contains all branch of seq.
	var travel func(*Seq, string)

	travel = func(n *Seq, indent string) {
		myPath := fmt.Sprintf("%s -> %s.%s", indent, n.DName, n.EName)
		if len(n.Cases) == 0 {
			if n.Next == nil {
				res = append(res, myPath)

				// end of the branch.
				return
			}

			// continue to the next state.
			travel(n.Next, myPath)
			return
		}
		myPath += " -> (select cases)"
		res = append(res, myPath)
		branchIndent := util.MkSpace(len(myPath))

		i := 1
		util.LoopStrMap(n.Cases, func(_ string, c *Seq) {
			travel(c, branchIndent+fmt.Sprintf("%d)", i))
			i++
		})
	}

	travel(s, indent)

	return res
}
