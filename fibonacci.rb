def fib(n)
  if n <= 1
    return n
  end

  fib(n - 1) + fib(n - 2)
end

fib(28)
