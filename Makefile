DIR := test
SRCS := $(wildcard $(DIR)/*/*.s)
TSTS := $(wildcard $(DIR)/*/*.go)
OBJS := $(patsubst $(DIR)/%.s, $(DIR)/%.bin,$(SRCS))
insttests: $(OBJS)
%.bin: %.s 
	riscv64-unknown-elf-gcc -o test/fib/fib.s -S test/fib/fib.c
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -march=rv64i -mabi=lp64 -o test/fib/fib test/fib/fib.s
	riscv64-unknown-elf-objcopy -O binary test/fib/fib test/fib/fib.bin

	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -march=rv64im -mabi=lp64 -o $*.out $<
	riscv64-unknown-elf-objcopy -O binary $*.out $@

add-addi.bin: test/add-addi.s
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -march=rv64i -mabi=lp64 -o test/add-addi test/add-addi.s
	riscv64-unknown-elf-objcopy -O binary test/add-addi test/add-addi.bin

testinst: 
	$(foreach var, $(TSTS),go run $(var);)

lint:
	golangci-lint run -v

clean:
	rm -f test/add-addi
	rm -f test/add-addi.bin
	rm -f test/fib/fib.bin
	rm -f test/fib/fib.s
	rm -f test/fib/fib
	rm -f $(DIR)/*/*.out
	rm -f $(DIR)/*/*.bin
