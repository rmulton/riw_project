package buildingIndexes

// BuildingIndex is an interface to implement to have an index that can be built
type BuildingIndex interface {
	AddDocToTerm(int, string)
	AddDocToIndex(int, string)
	GetDocCounter() int
}
