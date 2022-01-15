FROM alpine:latest
RUN mkdir -p /usr/bin/
ADD bin/kryptonite /usr/bin/
ADD bin/node-eth /usr/bin/
WORKDIR /usr/bin/
CMD [ "kryptonite" ]