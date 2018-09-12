package character

// Region defines an enumeration for the supported CAS region identifiers
type CreationState int

// Enumeration constants for CAS region identifiers
const (
	CreationStateStart CreationState = iota
	CreationStateEnd
)

var stateOrder = []CreationState{
	CreationStateStart,
	CreationStateEnd,
}

func (n CreationState) Next() CreationState {
	i := int(n)
	return stateOrder[i+1]
}
