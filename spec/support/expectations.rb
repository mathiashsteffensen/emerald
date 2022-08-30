module Expectations
    class Expectation
        attr_reader(:matcher)

        def initialize(matcher)
            @matcher = matcher
        end

        def matches?(expected)
            matcher.call(expected)
        end
    end

    class EqualityMatcher
        attr_reader(:actual)

        def initialize(actual)
            @actual = actual
        end

        def call(expected)
            expected == actual
        end
    end

    def eq(actual)
        Expectation.new(EqualityMatcher.new(actual))
    end
end
