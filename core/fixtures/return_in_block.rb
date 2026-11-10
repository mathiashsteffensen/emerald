$return_in_block_loaded = true
values = [1]
values.each { |x| return }
$return_in_block_loaded = :after_return
