package Structs

type Journaling struct {
	Operation [15]byte
	Path      [120]byte
	Content   [20]byte
	Date      [50]byte
}

func NewJournaling() Journaling {
	var journaling Journaling
	return journaling
}
