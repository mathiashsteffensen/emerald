x = 5
y = 7

LEVELS = [:fatal, :error, :warn, :info, :debug, :trace]

module BaseLogger
    attr_accessor(:level)

    100.times do
        LEVELS.each do |lvl|
            define_method(lvl) do |msg|
                if should_log?(lvl)
                    puts(lvl.to_s.upcase + " | " + msg)
                end
            end
        end
    end

    def should_log?(lvl)
        index_for_level(lvl) <= index_for_level(level)
    end

    def index_for_level(lvl)
        LEVELS.find_index { |other_level| other_level == lvl }
    end
end

class Logger
    include(BaseLogger)

    class << self
        include(BaseLogger)

        def instance
            @instance ||= new
        end
    end
end

Logger.instance.level = :info
Logger.level = :info

if x < y && true
    if true
        Logger.instance.info("Hello World!")

        Logger.instance.debug("Won't see this")
        Logger.instance.level = :debug
        Logger.instance.debug("Hello again")

        ["this", "is", "an", "array"].each { |msg| Logger.warn(msg) }
    end
else
    y
end
