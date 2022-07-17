FROM alpine:latest
ARG MIRROR=dl-cdn.alpinelinux.org
RUN sed -i "s/dl-cdn.alpinelinux.org/$MIRROR/g" /etc/apk/repositories
RUN apk add zsh-vcs git python3 fzf
WORKDIR /app
COPY . .
ENV DISABLE_OMZ_AUTOUPDATE=1 \
    USE_PLUGLOADER=1 \
    USE_RECOMMEND_PLUGINS=1 \
    USE_COMMON_ALIASES=1 \
    USE_LOCALIZE=1 \
    USE_RECOMMEND_THEME=1 \
    USE_SYNTAX_HIGHLIGHT=1 \
    USE_AUTO_SUGGEST=1 \
    USE_PACMAN_PKGFILE=0
CMD ./go
