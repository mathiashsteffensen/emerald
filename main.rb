x = 5
y = 7

class Logger
    def info(msg)
        puts("INFO ¦ " + msg)
    end
end

if x < y
    if true
        Logger.info("Hello World!")
    end
else
    y
end
