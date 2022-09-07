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
	        file = ARGV[0]
	        
	        if file
	        	puts "Running " + file
	        	require_relative "../../" + file
	        	puts "Done!"
	        else
	        	puts "Don't know how to run all specs yet :/"
	        end
	    end
	end
end

class Object
	include EMSpec::SpecHelpers
end
