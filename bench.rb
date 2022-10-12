def measure_time(name)
    start = Time.now.to_f
    yield
    puts("#{name} completed in #{Time.now.to_f - start}s")
end

def cache
    @cache ||= {}
end

def cached_fib_iterative
    cache[n] ||= fib_iterative
end

def fib_iterative(n)
    initial = [1,0]

    result = Range.new(0, n).reduce(initial) do |acc, w|
        newAccumulated = acc[1]
        nextAccumulated = acc[0]+acc[1]
        [newAccumulated, nextAccumulated]
    end

    result[0]
end

def bench_fib_iterative(n)
    measure_time "Iterative fibonacci #{n}" do
        cached_fib_iterative n
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

def string_add(n)
    n.times do
        "hello" + " " + "world" + "."
    end
end

def bench_string_add(n)
    measure_time "String add #{n * 3}" do
        string_add n
    end
end

5.times do
    # bench_fib_iterative(160_800)
    # bench_boolean_negate(1_000_000)
    bench_string_add 100_000
end

sleep 1
