start = Time.now.to_f

count = 1_000_000

count.times do
    "hello" == "hello"
end

puts("Compared " + count.to_s + " Strings in " + (Time.now.to_f - start).to_s + "s")

count.times do
    :hello == :hello
end

puts("Compared " + count.to_s + " Symbols in " + (Time.now.to_f - start).to_s + "s")
