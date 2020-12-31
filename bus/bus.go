package bus

//const debug bool = false

//Device for the bus need to implement load and store
type Device interface {
	// Methods signature with data types of the methods .
	Load(uint64, uint64) (uint64, error)
	Store(uint64, uint64, uint64) error
}
