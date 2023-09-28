module Response
    class << self
        def new(content_type, body)
            "HTTP/1.1 200 OK
Server: Emerald/#{Emerald.version}
Cache-Control: no-cache
Content-Type: #{content_type}
Content-Length: #{body.size}
Connection: Closed
#{body}
"
        end
    end
end