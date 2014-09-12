class Stream
  def initialize
    images = (1..1500).map{|i| "image number #{i}" }
    @string = StringIO.new(images.join("\n"))
  end

  def size
    @string.size
  end

  def read(length, outbuf)
    # sleep 0.1
    @string.read(30, outbuf)
  end
end

require "net/http"
require "uri"
require "stringio"

uri = URI.parse("http://localhost:7000/p/batch")
http = Net::HTTP.new(uri.host, uri.port)
request = Net::HTTP::Post.new(uri.request_uri)
request.content_type = "text/plain"

slow_stream = Stream.new
request.content_length = slow_stream.size
request.body_stream = slow_stream

p response = http.request(request)
puts response.body
