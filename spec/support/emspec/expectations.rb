module EMSpec
    module Expectations
        class RootExpectation
            attr_reader :expected

            def initialize(expected)
                @expected = expected
            end

            def to(matcher)
                if !matcher.matches?(expected)
                    matcher.error(expected)
                end
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
                puts "Expected " + expected.inspect + " to equal " + actual.inspect
            end
        end

        def expect(expected)
            Expectation.new(expected)
        end

        def eq(actual)
            EqualityMatcher.new(actual)
        end
    end
end
