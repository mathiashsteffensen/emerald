def measure_time(name)
    start = Time.now.to_f
    yield
    puts(name + " completed in " + (Time.now.to_f - start).to_s + "s")
end

def fib_recursive(n)
  if n <= 1
    return n
  end

  fib_recursive(n - 1) + fib_recursive(n - 2)
end

def fib_iterative(n)
    Range.new(0, n).inject([1,0]) { |acc, _| [acc[0], acc[0]+acc[1]] }[0]
end

def bench_fib_recursive(n)
    measure_time "Recursive fibonacci " + n.to_s do
        fib_recursive n
    end
end

def bench_fib_iterative(n)
    measure_time "Iterative fibonacci " + n.to_s do
        fib_iterative n
    end
end

bench_fib_iterative 168
