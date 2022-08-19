count = 100_000


str = "hello"

start = Time.now.to_f
count.times do
    str == str
end
puts("Compared " + count.to_s + " Strings in " + (Time.now.to_f - start).to_s + "s")


sym = :hello

start = Time.now.to_f
count.times do
    sym == sym
end
puts("Compared " + count.to_s + " Symbols in " + (Time.now.to_f - start).to_s + "s")
