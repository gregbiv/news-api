FROM golang:1.9-alpine

ENV APP_DIR $GOPATH/src/github.com/gregbiv/news-api
ENV MAS_COMMAND --version

# install Glide
RUN apk update \
    && apk add curl \
    && curl https://glide.sh/get | sh

# Install dependencies
RUN apk add --update alpine-sdk bzr libmagic file-dev \
    && rm -rf /var/cache/apk/*

# Setup file system
RUN mkdir -p ${APP_DIR} \
    && mkdir /var/log/news-api \
    && chmod 755 /var/log/news-api

COPY . ${APP_DIR}
WORKDIR ${APP_DIR}

RUN make deps
RUN make deps-dev

CMD CompileDaemon -build="make install" -command="news-api ${MAS_COMMAND}" -exclude-dir="pkg/assets/docs"
