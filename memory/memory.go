package memory

const debug bool = false

//MemorySize defines the  memory size for the cpu
const memorySize uint64 = 1024 * 1024 * 128

//Register holds all registers of the cpu
type memory struct {
	//Add memory here for simplicity
	//memory []uint8
	memory uint64
}

var mem memory

func (m *memory) Load() uint64        { return mem.memory }
func (m *memory) Store(binary uint64) { mem.memory = binary }

//Initialize the Memory
//func Initialize(binary []uint8) {
//	mem.memory = binary
//}
