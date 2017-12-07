package indexes

import (
	"fmt"
	"sort"
)

// ReversedIndex stores the index built by a parser.
// 
//		{
//				"toto": {
//					1234: 4,
//					23: 1,				// ReversedIndex structure
//					34: 10,
//					...
//				},
//		
//				...
//			}
type ReversedIndex map[string]map[int]int

func (index ReversedIndex) String() string {
	// Output variable
	var output string

	// Get all the keys and order them
	var keys []string
	for k := range index {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Append key-value pairs to the output
	for _, key := range keys {
		termDict := index[key]
		output += fmt.Sprintf("%s : %s\n", key, fmt.Sprint(termDict))
	}
		
	return output
}