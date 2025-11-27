package ledger

import (
	"fmt"
	"strconv"
)

type data_catalog struct {
	Grambank map[string] []*data_source
}

func new_data_catalog() *data_catalog {
	return &data_catalog {
		Grambank: make(map[string] []*data_source),
	}
}

func grammarize(whole string) []string {
	const gramsize = 3
	var grams []string

	if len(whole) < gramsize {
		return []string{ whole }
	}

	grams = make([]string, 0, len(whole) - gramsize + 1)

	for i := 0; i <= len(whole)-gramsize; i++ {
		gram := whole[i : i+gramsize]
		grams = append(grams, gram)
	}

	return grams
}

func (dc *data_catalog) AddSource(ds *data_source) {
	var gs []string
	var gram string

	if (ds.Name == "") {
		return
	}

	gs = grammarize(ds.Name)
	for _, gram = range gs {
		if (cap(dc.Grambank[gram]) == 0) {
			dc.Grambank[gram] = make([]*data_source, 0)
		}
		dc.Grambank[gram] = append(dc.Grambank[gram], ds) // TODO: this is where dups get messy
	}
}

// -------------- Search Functionality --------------- //
type SearchResult struct {
	Name	string		`json:"name"`
	Addr	string		`json:"address"`
	Etype	EntityType	`json:"entity_type"`
}

func (dc *data_catalog) SearchK(term string) []SearchResult {
	var (
		count			int
		gram			string
		dsrc			*data_source
		term_grammar	[]string
		search_heap		*search_heap
		search_counts	map[*data_source]int
		result			[]SearchResult
		search_pop		search_heap_item
		err				error
	)

	const top_k = 4

	term_grammar = grammarize(term)
	search_counts = make(map[*data_source]int)
	for _, gram = range term_grammar {
		for _, dsrc = range dc.Grambank[gram] {
			search_counts[dsrc] = search_counts[dsrc] + 1
		}
	}

	// if there are no search matches return early
	if len(search_counts) == 0 {
		return []SearchResult{}
	}

	search_heap = new_search_heap(len(search_counts))

	for dsrc, count = range search_counts {
		search_heap.Push(search_heap_item {
			Dsrc: dsrc,
			Count: count,
		})
	}
	
	result = make([]SearchResult, min(top_k, search_heap.count))
	for i := 0; i < top_k; i++ {
		search_pop, err = search_heap.Pop()
		if err != nil {
			break
		}
		result[i] = SearchResult{
			Name: search_pop.Dsrc.Name,
			Addr: strconv.FormatUint(uint64(search_pop.Dsrc.Addr), 10),
			Etype: search_pop.Dsrc.Etype,
		}
	}
	return result
}

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}


// -------------------------- Search Heap ---------------------------- //

type search_heap_item struct {
	Dsrc	*data_source
	Count	int
}

func (shi search_heap_item) String() string {
	return fmt.Sprintf("'%v' with count of %v", shi.Dsrc, shi.Count)
}

type search_heap struct {
	items []search_heap_item
	count int
}

func new_search_heap(capacity int) *search_heap {
	return &search_heap {
		items: make([]search_heap_item, capacity),
		count: 0,
	}
}

func (sh *search_heap) Push(item search_heap_item) error {
	var (
		crnt_index	int
		prnt_index	int
		temp		search_heap_item
	)

	if sh.count >= cap(sh.items) {
		return fmt.Errorf("failed to push, search heap is full")
	}

	// new root
	if sh.count == 0 {
		sh.items[0] = item
		sh.count += 1
		return nil
	}

	// place at leaf
	crnt_index = sh.count
	sh.items[crnt_index] = item
	sh.count += 1

	// fix up
	for {
		if crnt_index == 0 {
			break
		}

		prnt_index = (crnt_index - 1) / 2
		if sh.items[prnt_index].Count >= sh.items[crnt_index].Count {
			break // this means the heap is ordered
		}

		// swap
		temp = sh.items[crnt_index]
		sh.items[crnt_index] = sh.items[prnt_index]
		sh.items[prnt_index] = temp

		crnt_index = prnt_index
	}
	return nil
}

func (sh *search_heap) Pop() (search_heap_item, error) {
	var (
		ret 		search_heap_item
		chld_left	search_heap_item
		chld_rght	search_heap_item
		crnt_index	int
		rchld_cnt	int
		lchld_cnt	int
	)
	
	if sh.count == 0 {
		return search_heap_item{}, fmt.Errorf("failed to pop, search heap is empty.")
	}

	ret = sh.items[0]
	sh.count -= 1
	sh.items[0] = sh.items[sh.count]

	crnt_index = 0
	for {
		if crnt_index >= sh.count {
			break
		}
		
		if (crnt_index + 1) * 2 < sh.count {
			chld_rght = sh.items[(crnt_index + 1) * 2]
			rchld_cnt = chld_rght.Count
		} else {
			rchld_cnt = 0
		}

		if crnt_index * 2 + 1 < sh.count {
			chld_left = sh.items[crnt_index * 2 + 1]
			lchld_cnt = chld_left.Count
		} else {
			lchld_cnt = 0
		}

		if (rchld_cnt > lchld_cnt) && (sh.items[crnt_index].Count < rchld_cnt) {
			sh.items[(crnt_index + 1) * 2] = sh.items[crnt_index]
			sh.items[crnt_index] = chld_rght
			crnt_index = (crnt_index + 1) * 2
		} else if (rchld_cnt < lchld_cnt) && (sh.items[crnt_index].Count < lchld_cnt) {
			sh.items[crnt_index * 2 + 1] = sh.items[crnt_index]
			sh.items[crnt_index] = chld_left
			crnt_index = crnt_index * 2 + 1
		} else {
			break
		}
	}

	return ret, nil
}