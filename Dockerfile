FROM alpine:latest
ARG MIRROR=dl-cdn.alpinelinux.org
RUN sed -i "s/dl-cdn.alpinelinux.org/$MIRROR/g" /etc/apk/repositories
RUN apk add zsh-vcs git python3
WORKDIR /app
COPY . .
