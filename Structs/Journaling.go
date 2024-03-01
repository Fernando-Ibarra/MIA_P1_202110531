package Structs

type Journaling struct {
	Operation [20]byte
	Path      [150]byte
	Content   [10]byte
	Date      [150]byte
}
