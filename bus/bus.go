package bus

const debug bool = false

//MemoryBase set base addr for memroy, like QEMU vm
const MemoryBase uint64 = 0x8000_0000

//const MemoryBase uint64 = 0

//Device for the bus need to implement load and store
type Device interface {
	// Methods signature with data types of the methods .
	Load(uint64, uint64) (uint64, error)
	Store(uint64, uint64, uint64) error
}
