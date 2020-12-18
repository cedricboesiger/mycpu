main:
  addi x25, x0, 0xfe
  addi x26, x0, 0xff
  addi x27, x0, 2
  addi x28, x0, 42
  sb x27, 0(x25)
  sb x28, 0(x26)
  lb x30, 0(x25)
  lb x31, 0(x26)
