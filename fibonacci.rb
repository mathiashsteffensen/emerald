def measure_time(name)
    start = Time.now.to_f
    yield
    puts(name + " completed in " + (Time.now.to_f - start).to_s + "s")
end

def fib_iterative(n)
    Range.new(0, n).inject([1,0]) { |acc, w| [acc[1], acc[0]+acc[1]] }[0]
end

def bench_fib_iterative(n)
    measure_time("Iterative fibonacci " + n.to_s) do
        fib_iterative n
    end
end

5.times do
    bench_fib_iterative(160_800)
end

sleep 3
