module EMSpec
	class Context
		attr_reader :name, :parent
		attr_accessor :child

		def initialize(name, parent, child)
			@name = name
			@parent = parent
			@child = child
		end

		module Helpers
			def context(name)
				@current_context = EMSpec::Context.new(name, @current_context, nil)

				with_child { yield }
			end

			def describe(name)
				context(name) { yield }
			end

			def current_context
				@current_context ||= EMSpec::Context.new("EMSpec", nil, nil)
			end

			def with_child
				@current_context = current_context.child
				yield
				@current_context = current_context.parent
			end
		end
	end
end
