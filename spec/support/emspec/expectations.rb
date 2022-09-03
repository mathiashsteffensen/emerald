module EMSpec
    module Expectations
        class Expectation
            attr_reader :expected

            def initialize(expected)
                @expected = expected
            end

            def to(matcher)
                if !matcher.matches?(expected)
                    log_failure(matcher.error(expected))
                    nil
                end
            end

            def log_failure(msg)
                puts "Spec failed"
                puts "  " + current_context.name
                puts "      " + msg
            end
        end

        class EqualityMatcher
            attr_reader :actual

            def initialize(actual)
                @actual = actual
            end

            def matches?(expected)
                expected == actual
            end

            def error(expected)
                "Expected " + expected.inspect + " to equal " + actual.inspect
            end
        end

        def expect(expected)
            EMSpec::Expectations::Expectation.new(expected)
        end

        def eq(actual)
            EMSpec::Expectations::EqualityMatcher.new(actual)
        end
    end
end
