FROM alpine:latest
ARG MIRROR=dl-cdn.alpinelinux.org
RUN sed -i "s/dl-cdn.alpinelinux.org/$MIRROR/g" /etc/apk/repositories
RUN apk update
WORKDIR /app
COPY . .
RUN yes | ./setup
CMD /bin/zsh
