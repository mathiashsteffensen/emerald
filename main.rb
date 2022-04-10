x = 5
y = 7

LEVELS = [:fatal, :error, :warn, :info, :debug, :trace]

class Logger
    LEVELS.each do |lvl|
        define_method(lvl.to_s) do |msg|
            if should_log(lvl)
                puts(lvl.to_s.upcase + " | " + msg)
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
        @instances = 0

        def instance
            if @instance == nil
                @instances = @instances + 1
                @instance = new
            end

            @instance
        end

        def log_instances
            instance.info("You have created " + @instances.to_s + " instances of Logger")
        end
    end
end

if x < y
    if true
        Logger.instance.info("Hello World!")

        Logger.instance.debug("debug msg")

        ["this", "is", "an", "array"].each { |msg| Logger.instance.warn(msg) }

        Logger.log_instances
    end
else
    y
end
