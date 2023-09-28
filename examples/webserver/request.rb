class Request
    attr_reader :method, :path, :protocol  

    def initialize(raw_request)
        parts = raw_request.split(" ")
        @method = parts[0]
        @path = parts[1]
        @protocol = parts[2]
    end
end
