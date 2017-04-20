FROM alpine:latest

ENV app /app

RUN mkdir -p $app
WORKDIR $app

RUN apk update
RUN apk add imagemagick curl vim
RUN curl -o $app/images http://us-east.manta.joyent.com/wanelo/public/images/bin/images-linux-1.13.11 && chmod +x $app/images

# Add Containerpilot and set its configuration
ENV CONTAINERPILOT_VERSION 2.7.0
ENV CONTAINERPILOT file:///etc/containerpilot.json

RUN export CONTAINERPILOT_CHECKSUM=687f7d83e031be7f497ffa94b234251270aee75b \
    && export archive=containerpilot-${CONTAINERPILOT_VERSION}.tar.gz \
    && curl -Lso /tmp/${archive} \
         "https://github.com/joyent/containerpilot/releases/download/${CONTAINERPILOT_VERSION}/${archive}" \
    && echo "${CONTAINERPILOT_CHECKSUM}  /tmp/${archive}" | sha1sum -c \
    && tar zxf /tmp/${archive} -C /usr/local/bin \
    && rm /tmp/${archive}

# configuration files and bootstrap scripts
COPY etc/containerpilot.json /etc/
COPY bin/health.sh /usr/local/bin/health.sh

CMD /usr/local/bin/containerpilot $app/images --outputs='' --listen 0.0.0.0 server

EXPOSE 7000 7002

