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
        file = ARGV[3]

        if file
          puts "Running " + file
          require_relative "../../" + file
        else
          Dir.glob("spec/**/*_spec.rb").each do |f|
            puts "Running " + f
            require_relative "../../" + f
          end
        end

        puts "Done!"
	    end
	end
end

class Object
	include EMSpec::SpecHelpers
end
