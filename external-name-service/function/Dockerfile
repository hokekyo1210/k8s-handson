# dev, builder
FROM golang:1.14
WORKDIR /work/app
ENV GO111MODULE=on

RUN curl -fLo /usr/local/bin/air https://github.com/cosmtrek/air/raw/v1.12.1/bin/linux/air \
    && chmod +x /usr/local/bin/air

COPY . /work/app

ENTRYPOINT ["go", "run", "/work/app/cmd/main.go"]