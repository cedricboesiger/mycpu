use std::env;
use std::io;
use std::fs::File;
use std::io::prelude::*;

// Default memory size
pub const MEMORY_SIZE: u64 = 1024 * 1024 * 128; //(128MiB)

// Cpu contains registers,  program counter and memory
struct Cpu {
    // 32 64bit registers
    regs: [u64; 32],
    // One program counter
    pc: u64,
    // Memory for simplicity
    memory: Vec<u8>,
}

impl Cpu {

    fn new(binary: Vec<u8>) -> Self {
        // Initialize registers
        let mut regs = [0;32];
        // Set register x02 as sp
        regs[2] = MEMORY_SIZE;

        Self{
            regs,
            pc: 0,
            memory: binary,
        }
    }

    /// Print values in all registers (x0-x31).
    pub fn dump_registers(&self) {
        let mut output = String::from("");
        let abi = [
            "zero", " ra ", " sp ", " gp ", " tp ", " t0 ", " t1 ", " t2 ", " s0 ", " s1 ", " a0 ",
            " a1 ", " a2 ", " a3 ", " a4 ", " a5 ", " a6 ", " a7 ", " s2 ", " s3 ", " s4 ", " s5 ",
            " s6 ", " s7 ", " s8 ", " s9 ", " s10", " s11", " t3 ", " t4 ", " t5 ", " t6 ",
        ];
        for i in (0..32).step_by(4) {
            output = format!(
                "{}\n{}",
                output,
                format!(
                    "x{:02}({})= {:<#18x} x{:02}({})= {:<#18x} x{:02}({})= {:<#18x} x{:02}({})= {:<#18x}",
                    i,
                    abi[i],
                    self.regs[i],
                    i + 1,
                    abi[i + 1],
                    self.regs[i + 1],
                    i + 2,
                    abi[i + 2],
                    self.regs[i + 2],
                    i + 3,
                    abi[i + 3],
                    self.regs[i + 3],
                )
            );
        }
        println!("{}", output);
    }

    //Get instruction from memory
    fn fetch(&self) -> u32 {
        let index = self.pc as usize;
        return (self.memory[index] as u32)
            | ((self.memory[index + 1] as u32) << 8)
            | ((self.memory[index + 2] as u32) << 16)
            | ((self.memory[index + 3] as u32) << 24);
    }

    //Execute instruction
    fn execute(&mut self, inst: u32) {
        // Set register x00 as hardwired to 0
        self.regs[0] = 0;

        let opcode = inst & 0x0000007f;
        let rd = ((inst & 0x00000f80) >> 7) as usize;
        let rs1 = ((inst & 0x000f8000) >> 15) as usize;
        let rs2 = ((inst & 0x01f00000) >> 20) as usize;

	match opcode {
            0x13 => {
                // addi
                let imm = ((inst & 0xfff00000) as i32 as i64 >> 20) as u64;
                self.regs[rd] = self.regs[rs1].wrapping_add(imm);
            }
            0x33 => {
                // add
                self.regs[rd] = self.regs[rs1].wrapping_add(self.regs[rs2]);
            }
            _ => {
                dbg!(format!("not implemented yet: opcode {:#x}", opcode));
            }
        }
    }
}

fn main() -> io::Result<()> {
    
    let args: Vec<String> = env::args().collect();

    if args.len() != 2 {
        panic!("Usage: mycpu <binary-filename>");
    }

    let mut file = File::open(&args[1])?;
    let mut binary = Vec::new();
    file.read_to_end(&mut binary)?;

    let mut cpu = Cpu::new(binary);

    while cpu.pc < cpu.memory.len() as u64 {

        // 1. Fetch
        let inst = cpu.fetch();

        // 2. Add 4 to the program counter
        cpu.pc += 4;

        // 3. Decode
        // 4. Execute
        cpu.execute(inst);
    }

    cpu.dump_registers();

    Ok(())

}


