module EMSpec
	class << self
		def _current_context
			$current_context ||= EMSpec::Context.new("EMSpec", nil, nil)
		end
	end

	class Context
		attr_reader :parent, :name, :level

		def initialize(name, parent, level)
			@name = name
			@level = level
		end

		module Helpers
			def context(name)
				$current_context = EMSpec::Context.new(name, current_context, current_context.level + 1)

				with_context_reset { yield; nil }
			end

			def describe(name)
				context(name) { yield }
			end

			def it(name)
				context(name) { yield }
			rescue RuntimeError
				log_failure($!.to_s)
			end

			def current_context
				EMSpec._current_context
			end

			def with_context_reset
				$current_context = current_context.child
				yield
				$current_context = current_context.parent
			end

			def log_failure(msg)
                puts "Spec failed"
                puts "  " + current_context.name
                puts "      " + msg
            end
		end
	end
end
