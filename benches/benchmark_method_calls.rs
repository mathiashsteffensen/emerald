use criterion::{criterion_group, criterion_main, Criterion};
use emerald::kernel;

fn fibonacci() {
    kernel::execute(
        "bench.rb".to_string(),
        "def fib(n)
          if n <= 1
            return n
          end

          fib(n - 1) + fib(n - 2)
        end

        fib(10)"
            .to_string(),
    )
    .unwrap();
}

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("fib", |b| b.iter(|| fibonacci()));
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
