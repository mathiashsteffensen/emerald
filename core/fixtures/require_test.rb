def rescue_error_to_see_if_doing_so_in_require_fucks_up_stack
    raise "Shit going down"
rescue StandardError
    puts("Shit went down my guy")
end

rescue_error_to_see_if_doing_so_in_require_fucks_up_stack
require_relative "require_test_2"
