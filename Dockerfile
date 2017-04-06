FROM golang:1.7

ENV app /app

RUN mkdir -p $app
WORKDIR $app
ADD . $app

RUN apt-get update
RUN apt-get install imagemagick -y
RUN make deps
RUN make build

CMD $app/bin/linux/images --outputs='' server

EXPOSE 7000 7002

