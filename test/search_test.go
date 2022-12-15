package test

import (
	"fmt"
	"go_tests/search"
	"testing"
)

var searchLimit *search.SearcherLimit

func init() {
	var err error
	searchLimit, err = search.NewSearcherLimit(&SimpleStruct{})
	if err != nil {
		fmt.Println(err)
		panic("search.NewSearcherLimit error")
	}
	fmt.Printf("searchLimit: %+v\n", *searchLimit)
}

func BenchmarkSearch(b *testing.B) {
	b.ReportAllocs()
	datas := []*SimpleStruct{
		{A: 1, B: 10, Str: "wzyao1"},
		{A: 2, B: 20, Str: "wzyao2"},
		{A: 3, B: 30, Str: "wzyao3"},
		{A: 4, B: 40, Str: "wzyao4"},
		{A: 5, B: 50, Str: "wzyao5"},
	}
	for i := 0; i < b.N; i++ {
		searchs := []*search.Searcher{
			{Field: "a", SearchOperator: search.SEARCH_OPERATOR_LESS_EQUAL, Value: "3"},
			{Field: "b", SearchOperator: search.SEARCH_OPERATOR_LESS_EQUAL, Value: "20"},
			{Field: "str", SearchOperator: search.SEARCH_OPERATOR_CONTAIN_OR, Value: "wzyao"},
		}
		if err := searchLimit.ValidCheck(searchs); err != nil {
			b.Fatalf("searchLimit.ValidCheck: %s", err.Error())
		}
		var datasIn []interface{}
		for _, v := range datas {
			datasIn = append(datasIn, v)
		}
		var err error
		for _, s := range searchs {
			datasIn, err = s.Filter(searchLimit, datasIn)
			if err != nil {
				fmt.Println("filter err:", err.Error())
				return
			}
		}
	}
}
