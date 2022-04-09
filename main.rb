x = 5
y = 7

LEVELS = [:fatal, :error, :warn, :info, :debug, :trace]

class Logger
    LEVELS.each do |lvl|
        define_method(lvl.to_s) do |msg|
            puts(lvl.to_s.upcase + " | " + msg)
        end
    end

    class << self
        def debug(msg)
            if index_for_level(:debug) <= index_for_level(level)
                puts("DEBUG ¦ " + msg)
            end
        end

        def info(msg)
            if index_for_level(:info) <= index_for_level(level)
                puts("INFO ¦ " + msg)
            end
        end

        def level
            :info
        end

        def index_for_level(lvl)
            LEVELS.find_index { |other_level| other_level == lvl }
        end
    end
end

instance = Logger.new

if x < y
    if true
        Logger.info("Hello World!")
        Logger.debug("Hello World!")

        instance.debug("debug msg")

        ["this", "is", "an", "array"].each { |msg| instance.warn(msg) }
    end
else
    y
end
