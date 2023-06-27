package meta

type UpdateStrategy int8

func (u UpdateStrategy) String() string {
	
	return ""
}

const (
	AppendChanges UpdateStrategy = iota
	AppendAll
	ReplaceChanges
	ReplaceAll
)
