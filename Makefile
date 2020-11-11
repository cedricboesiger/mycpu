add-addi.bin: test/add-addi.s
	riscv64-unknown-elf-gcc -Wl,-Ttext=0x0 -nostdlib -o test/add-addi test/add-addi.s
	riscv64-unknown-elf-objcopy -O binary test/add-addi test/add-addi.bin

clean:
	rm -f test/add-addi
	rm -f test/add-addi.bin
