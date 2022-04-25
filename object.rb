start = Time.now.to_f

count = 1_000_000

count.times do
    Object.new
end

puts("Created " + count.to_s + " Objects in " + (Time.now.to_f - start).to_s + "s")
