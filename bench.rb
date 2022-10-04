def measure_time(name)
    start = Time.now.to_f
    yield
    puts("#{name} completed in #{Time.now.to_f - start}s")
end

def cache
    @cache ||= {}
end

def fib_iterative(n)
    cache[n] ||= Range.new(0, n).reduce([1,0]) { |acc, w| [acc[1], acc[0]+acc[1]] }[0]
end

def bench_fib_iterative(n)
    measure_time "Iterative fibonacci #{n}" do
        fib_iterative n
    end
end

def boolean_negate(n)
    n.times do
        !nil && !false

        !"boop" || !2 || !/re/
    end
end

def bench_boolean_negate(n)
    measure_time "Boolean negate #{n * 5}" do
        boolean_negate n
    end
end

5.times do
    # bench_fib_iterative(160_800)
    bench_boolean_negate(1_000_000)
end

sleep 1
