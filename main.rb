x = 5
y = 7

class Logger
    def debug(msg)
        puts("DEBUG ¦ " + msg)
    end

    class << self
        def info(msg)
            puts("INFO ¦ " + msg)
        end
    end
end

instance = Logger.new

if x < y
    if true
        Logger.info("Hello World!")
        instance.debug(" from debug land")
    end
else
    y
end
