require "net/http"
require "uri"
require "stringio"

uri = URI.parse("http://localhost:7000/p/batch")
# uri = URI.parse("http://localhost:3333/p/batch")
http = Net::HTTP.new(uri.host, uri.port)
request = Net::HTTP::Post.new(uri.request_uri)

class Stream
  def initialize
    @string = StringIO.new("hey there buddy\n"*20 + "\n\n")
  end

  def size
    @string.size
  end

  def read(length, outbuf)
    sleep 0.1
    @string.read(30, outbuf)
  end
end

slow_stream = Stream.new
request.content_length = slow_stream.size
request.content_type = "text/plain"
request.body_stream = slow_stream
p response = http.request(request)
