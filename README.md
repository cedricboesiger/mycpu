# rvsim a Risc-V simulator
a riscv emulator implemented in go.
it is inspired by the rvemu project, written in rust: [rvemu from scratch in 10 steps](https://book.rvemu.app/)

# v0.1.0, Step 01 reimplementd in Rust
Implemented STEP 1, it decodes and executes `add` and `addi` according to rv32i

* Create an assembler programm using add and addi, `add-addi.s`:
```
main:
  addi x29, x0, 2
  addi x30, x0, 40
  add x31, x30, x29
```
* then execute:
```
make clean && make
cargo run test/add-addi.bin
```
# v0.1.1, Step 01 implemented in go 
Same as v0.1.1 but implemented in go. Instead of cargo just run:
```
go run hart.go -f test/add-addi.bin
```
# v0.1.2, Step 02 implemented in go
```
go run hart.go -f test/fib/fib.bin
```
