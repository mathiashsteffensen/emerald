module Expectations
    class Expectation
        def initialize(matcher)
            @matcher = matcher
        end
    end

    def eq; end
end
