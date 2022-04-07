x = 5
y = 7

class Logger
    class << self
        def info(msg)
            puts("INFO ¦ " + msg)
        end
    end
end

if x < y
    if true
        Logger.info("Hello World!")
    end
else
    y
end
