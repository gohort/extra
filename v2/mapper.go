package extra

// Mapper represents the functionality of an extra Map. These functions don't
// have to be used but it can help for processing times if you embed an extra.Map
// into a struct.
type Mapper interface {
	// GetExtraMap returns a reference to a raw *Map object.
	GetExtraMap() *Map
	// RemoveExtraKey deletes the key provided.
	RemoveExtraKey(string)
	// GetExtraField returns a value and boolean if the key was found in this mapper.
	GetExtraField(string) (any, bool)
}
