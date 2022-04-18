x = 5
y = 7

LEVELS = [:fatal, :error, :warn, :info, :debug, :trace]

class Logger
    100.times do
        LEVELS.each do |lvl|
            define_method(lvl) do |msg|
                if should_log(lvl)
                    puts(lvl.to_s.upcase + " | " + msg)
                end
            end
        end
    end

    def should_log(lvl)
        index_for_level(lvl) <= index_for_level(current_level)
    end

    def index_for_level(lvl)
        LEVELS.find_index { |other_level| other_level == lvl }
    end

    def current_level
        :info
    end

    class << self
        def instance
            new
        end
    end
end

if x < y && true
    if true
        Logger.instance.info("Hello World!")

        Logger.instance.debug("debug msg")

        ["this", "is", "an", "array"].each { |msg| Logger.instance.warn(msg) }
    end
else
    y
end
