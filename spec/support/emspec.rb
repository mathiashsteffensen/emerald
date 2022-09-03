require_relative "emspec/context"
require_relative "emspec/expectations"

module EMSpec
	module SpecHelpers
		include Context::Helpers
		include Expectations
	end

	class << self
		include Context::Helpers

	    def run
	        file = ARGV[2] 
	        
	        if file
	        	puts "Running " file
	        	require_relative "../../" + file
	        end
	    end
	end
end

class Object
	include EMSpec::SpecHelpers
end
