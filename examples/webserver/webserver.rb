require_relative "request"
require_relative "response"

class Webserver
    attr_reader :server, :port

    def initialize(port)
        @port = port
    end
    
    def init
        @server = TCPServer.new "localhost", port
        self
    end
    
    def serve
        while true
            socket = server.accept
            req = read(socket)
        
            measure_time(req) do
                write(socket)

                socket.close
            end
        end
    end
    
    private
    
    def read(socket)
        Request.new(socket.gets)
    end
    
    def write(socket)
        socket.write(Response.new("text/html", get_page))
    rescue StandardError
        puts "Error writing to TCPSocket: #{$!}"
    end
    
    def get_page
        IO.read("./index.html")
    end
    
    def measure_time(req)
        start = Time.now.to_f
        
        yield
        
        elapsed = (Time.now.to_f-start).round(4)
        
        puts "#{req.method} #{req.path} - Processed request in #{elapsed}ms"
    end
end

Webserver
    .new(2000)
    .init
    .serve
