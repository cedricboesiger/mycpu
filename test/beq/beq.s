main:
  addi x31, x0, 0
  addi x28, x0, 2
  addi x29, x0, 2
  beq x28, x29, true
  jal x0, end
true:
  add x31, x29, x28
  jal x0, end
end:
  addi x0, x0, 0
  
