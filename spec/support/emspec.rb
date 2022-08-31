require_relative "emspec/expectations"

module EMSpec
	module SpecHelpers
		include Expectations
	end

	class << self
	    def run
	        puts "Running all specs"
	    end
	end
end

include EMSpec::SpecHelpers
