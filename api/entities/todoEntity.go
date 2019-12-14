package entities

type EStatus int

const (
	TODO       EStatus = 1
	INPROGRESS EStatus = 2
	DONE       EStatus = 3
)

type ToDo struct {
	Hash        string
	Title       string  `bson:"title"`
	Description string  `bson:"description"`
	Status      EStatus `bson:"status"`
}
