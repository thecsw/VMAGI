import sys

def fib(n):
    if n == 0 or n == 1:
        return 1
    return fib(n-1) + fib(n-2)

print(fib(int(sys.argv[1])))
