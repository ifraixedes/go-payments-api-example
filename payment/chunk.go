package payment

// Chunk is the type to indicate how many items of a collection to get and from
// which offset.
type Chunk struct {
	Limit  uint32
	Offset uint64
}
