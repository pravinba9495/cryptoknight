FROM alpine:latest
RUN mkdir -p /usr/bin/
ADD bin/kryptonite /usr/bin/
WORKDIR /usr/bin/
CMD [ "kryptonite" ]