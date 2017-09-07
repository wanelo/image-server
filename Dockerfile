FROM alpine:latest

ENV app /app

RUN mkdir -p $app
WORKDIR $app

RUN apk update
RUN apk add imagemagick curl
RUN curl -o $app/images http://us-east.manta.joyent.com/wanelo/public/images/bin/images-linux-1.13.11 && chmod +x $app/images

CMD $app/images --outputs='' --listen 0.0.0.0 server

EXPOSE 7000 7002

