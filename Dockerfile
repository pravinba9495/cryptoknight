FROM alpine:latest
RUN mkdir -p /usr/bin/
ADD kryptonite /usr/bin/
WORKDIR /usr/bin/
CMD [ "kryptonite" ]