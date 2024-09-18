package gojen_test

import (
	"fmt"
	"testing"

	"github.com/cirius-go/gojen"
)

func TestSequence(t *testing.T) {
	s := gojen.
		NewSeq("service", "init").
		Select("service", []string{"crud", "singleMethod"}, func(ss gojen.SeqSwitcher) {
			ss.When("crud", func(c *gojen.Seq) *gojen.Seq {
				return c.
					Append("dto", "init", "crud").
					Append("model", "init").
					Append("repo", "init").
					Append("api", "init", "crud")
			})

			ss.When("singleMethod", func(c *gojen.Seq) *gojen.Seq {
				return c.
					Append("dto", "init", "singleMethod").
					Append("api", "init").Select("api", []string{"get", "post", "put", "patch", "delete"}, nil)
			})
		})

	fmt.Println(s)

	// -> service.init -> (select cases)
	//                                  1) -> service.crud -> dto.init -> dto.crud -> model.init -> repo.init -> api.init -> api.crud
	//                                  2) -> service.singleMethod -> dto.init -> dto.singleMethod -> api.init -> (select cases)
	//                                                                                                                          1) -> api.delete
	//                                                                                                                          2) -> api.get
	//                                                                                                                          3) -> api.patch
	//                                                                                                                          4) -> api.post
	//                                                                                                                          5) -> api.put
}
