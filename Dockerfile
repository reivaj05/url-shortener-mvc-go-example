FROM alpine:3.3

ENV PROJECT_PATH=/go/src/github.com/reivaj05/url_shortener
ENV GOPATH=/go

ADD . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN apk -U add make gcc g++ icu-dev ncurses-dev git bash go curl \
    && cd $PROJECT_PATH \
    && make \
    && apk del make icu-dev ncurses-dev git go curl

CMD $PROJECT_PATH/url_shortener

EXPOSE 8000
