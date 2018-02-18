package building

type BuildingIndex interface {
	AddDocToTerm(int, string)
	AddDocToIndex(int, string)
	GetDocCounter() int
}