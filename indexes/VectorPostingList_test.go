package indexes

import (
	"math"
	"testing"
	"reflect"
)

var somePostingLists = map[string]PostingList {
	"blabla": PostingList {
		1: 2.5,
		2: 3.45,
		3: 6.934,
		342: 34.454,
	},
	"bleble": PostingList {
		1: 32.,
		3: 423.34,
	},
	"blublu": PostingList {
		1: 2.,
	},
}

var expectedMergedToVectorPostingList = VectorPostingList {
	1: map[string]float64 {
		"blabla": 2.5,
		"bleble": 32.,
		"blublu": 2.,
	},
	2: map[string]float64 {
		"blabla": 3.45,
	},
	3: map[string]float64 {
		"blabla": 6.934,
		"bleble": 423.34,
	},
	342: map[string]float64 {
		"blabla": 34.454,
	},
}

var reqVec = map[string]float64 {
	"blabla": 1.,
	"bleble": 1.,
	"blublu": 1.,
}

var expectedDocAngleScores = map[int]float64 {
    1: math.Acos((2.5 + 32. + 2.)/(math.Sqrt(math.Pow(2.5, 2.)+math.Pow(32.,2.)+math.Pow(2.,2))+math.Sqrt(3.))),
	2: math.Acos((3.45)/(math.Sqrt(math.Pow(3.45, 2.))+math.Sqrt(3.))),
	3: math.Acos((6.934 + 423.34)/(math.Sqrt(math.Pow(6.934, 2.)+math.Pow(423.34, 2.))+math.Sqrt(3.))),
	342: math.Acos((34.454)/(math.Sqrt(math.Pow(34.454,2.))+math.Sqrt(3.))),
} 

func TestMergeToVector(t *testing.T) {
	mergedToVectorPostingList := MergeToVector(somePostingLists)
	if !reflect.DeepEqual(mergedToVectorPostingList, expectedMergedToVectorPostingList) {
		t.Errorf("Merged to vector posting list should be %v, not %v", expectedMergedToVectorPostingList, mergedToVectorPostingList)
	}
}

func TestToAngleTo(t *testing.T) {
	vecPostingList := MergeToVector(somePostingLists)
	docsAngleScore := vecPostingList.ToAnglesTo(reqVec)
	if !reflect.DeepEqual(docsAngleScore, expectedDocAngleScores) {
		t.Errorf("Angles scores should be %v, not %v", expectedDocAngleScores, docsAngleScore)
	}
	
}
// test operations like angles
// test weird values on angles etc