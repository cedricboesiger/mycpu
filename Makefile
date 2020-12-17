DIR := test/e2e
SRCS := $(wildcard $(DIR)/*.s)
OBJS := $(patsubst $(DIR)/%.s, $(DIR)/%.bin,$(SRCS))
e2e: $(OBJS)
%.bin: %.s 
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -o $*.out $<
	riscv64-unknown-elf-objcopy -O binary $*.out $@

add-addi.bin: test/add-addi.s
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -o test/add-addi test/add-addi.s
	riscv64-unknown-elf-objcopy -O binary test/add-addi test/add-addi.bin

clean:
	rm -f test/add-addi
	rm -f test/add-addi.bin
	rm -f $(DIR)/*.out
	rm -f $(DIR)/*.bin
