DIR := test
SRCS := $(wildcard $(DIR)/*/*.s)
TSTS := $(wildcard $(DIR)/*/*.go)
OBJS := $(patsubst $(DIR)/%.s, $(DIR)/%.bin,$(SRCS))
insttests: $(OBJS)
%.bin: %.s 
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -o $*.out $<
	riscv64-unknown-elf-objcopy -O binary $*.out $@

add-addi.bin: test/add-addi.s
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -o test/add-addi test/add-addi.s
	riscv64-unknown-elf-objcopy -O binary test/add-addi test/add-addi.bin

testinst: 
	$(foreach var, $(TSTS),go run $(var);)

clean:
	rm -f test/add-addi
	rm -f test/add-addi.bin
	rm -f $(DIR)/*/*.out
	rm -f $(DIR)/*/*.bin
