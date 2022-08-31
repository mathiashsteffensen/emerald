module EMSpec
    module Expectations
        class RootExpectation
            attr_reader :expected, :expectation

            def initialize(expected)
                @expected = expected
            end

            def to(expectation)

            end
        end

        class Expectation
            attr_reader :matcher

            def initialize(matcher)
                @matcher = matcher
            end

            def matches?(expected)
                matcher.call(expected)
            end

            def error(expected)
                matcher.error(expected)
            end
        end

        class EqualityMatcher
            attr_reader :actual

            def initialize(actual)
                @actual = actual
            end

            def call(expected)
                expected == actual
            end

            def error(expected)
                puts "Expected " + expected.inspect + " to equal " + actual.inspect
            end
        end

        def eq(actual)
            Expectation.new(EqualityMatcher.new(actual))
        end
    end
end
