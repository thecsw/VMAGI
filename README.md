![VMAGI](./VMAGI.png)

# VMAGI

Welcome! `VMAGI` is a small emulator/interpreter my friend 
[Matthew](https://github.com/matthewsanetra) 
and I
challenged each other to build in 24 hours. This includes both the implementation of the interpreter
and creating your own ISA/3ac/IR for it that it will run on top of. 

If you want to see Matthew's implementation, go to [his repository](https://github.com/matthewsanetra/sandy_isa).

## Basics

The goal within 24 hours was to write an interpreter that can reliably run a recursive version
of fibonacci sequence, such that for any natural `n`, `fib(n)` returns the `n`th element of the
sequence. 

This was an interesting challenge, as writing the machine itself, with all the instructions and 
logic around it was pretty simple. This includes required stuff like labels, jumps, conditionals,
etc. The idea was to make it a workable interpreter for whatever you write. 

Matthew and I both have the same direction, in terms of what our interpreters should achieve. Besides
being able to run recursive fibonacci, it should have a command/opcode `halt` that will immediately
halt the interpreter and return the value that was passed to `halt` with either registers or some
other way.

In the end of the competition, which is right now as I'm writing it, `VMAGI` has a total of 25 opcodes.
You see, I have a very classical Computer Science training, so I think in terms of memory layout, IO-bound
operations, and etc. For more details, see `isa.go` for the list of instructions and how the look like.

## Example

Here is an example of the recursive fibonacci
```
        addi #1, 30, #1   -- 0 depth is 30
        push #1           -- 1 push depth to parameters
        call @fib!        -- 2 call fib(n)
        pop #1            -- 3 pop the return
        halt #1           -- 4 exit with the return

fib!:   pop #1            -- 5 pop the n param
        jmpz #1, @base0   -- 6 if n==0, jump to @base0
        eqi  #1, 1, #2    -- 7 set #2 = #1 == 1
        jmpt #2, @base1   -- 8 if #1 == 1, jump to @base1
        subi #1, 1, #1    -- 9 set #1 = n - 1
        push #1           -- 10 push n-1 to stack
        call @fib!        -- 11 call fib(n-1)
        pop #2            -- 12 pop result of fib(n-1) to #2
        subi #1, 1, #1    -- 13 set #1 = (n - 1) - 1 = n - 2
        push #1           -- 14 push n-2 to stack
        call @fib!        -- 15 call fib(n-2)
        pop #1            -- 16 pop result of fib(n-2) to #1
        add #1, #2, #1    -- 17 set #1 = #1 + #2 = fib(n-2) + fib(n-1)
        push #1           -- 18 push the fib result to stack
        ret               -- 19 return
        
base0:  pushi 1           -- 20 push 1 to stack
        ret               -- 21 return
        
base1:  pushi 1           -- 22 push 1 to stack
        ret               -- 23 return
```

To which `> VMAGI ./examples/fib.vmagi` outputs
```sh
VMAGI stopped execution with 1346269
```

## My idea

I tried to write `VMAGI` in such a way that its reasonably performant and most of all, dead simple. Simple
as in how it works internally, the way data gets shuffled around, and the cost of maintenance/understanding
the code. To build a simple stack machine is not hard, almost trivial in some places, as all you do is run
the most simple instructions, such as `add` or `push` and put the data somewhere, where it needs to go. 

One of the biggest mistakes you can make is usually just mistyping something or getting some logic nuances
nudged in a different directions. No exception here. I would say the biggest difficulty for this challenge
was the optimization stage and making everything faster, but also not sacrificing the code quality. Also,
because this is Go, I also tried to write everything using *just* the standard library, so no dependencies.

I just wanted to have fun and watch Kill la Kill for the 100th time. The basic version that just runs the code,
as taking in the assembly input, parsing it, and running it, all in all took about 4 hours. The next ~8 hours
were spent on cleaning up the code, improving the logic, tightnening up some holes, and coming up with some
optimizations. 

## How it works

`VMAGI` works in a direct, single-pipeline, sequential way. The only input `VMAGI` requires is the filename
with the `.vmagi` assembly code you want to be interpreted. `VMAGI` will read the file and take care of the
parsing and executing together. You don't need an assembler to pass in the bytecode. 

Parsing is fairly simple. See `parser.go` for that. What `VMAGI` does is read the given file line-by-line and
run a couple of regex expressions to figure out the opcode/command that is invoked with all its operands. Of 
course, you could do it better using a lexer, as Matthew did, like `flex` and use something ala `bison` to 
write down the grammar for our new language. I have done it many times, so I just wanted to have some fun.

After we parsed out the language (label lookups get a deferred parsing), we will have a list of all valid
instructions in our language (See `models.go`). After that, we will be starting on executing each instruction,
by moving the `PC` (Program Counter) pointer where it has to go. For more gory details on the execution logic
and code, please go and see `execute.go`

## Performance

Definitely one of the hardest and most annoying parts of writing `VMAGI` was improving performance on deeply
nested recursive calls. Recall that we had to make an emulator that recursively calculates any member of the
Fibonacci sequence.
