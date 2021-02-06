FROM golang:alpine as base

RUN mkdir /lastfm-collage
WORKDIR /lastfm-collage
COPY go.mod .
COPY go.sum .

RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
FROM base as pre-build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/lastfm-collage
FROM scratch as production
COPY --from=pre-build /go/bin/lastfm-collage /
COPY --from=base /lastfm-collage/static /static

ENTRYPOINT ["/lastfm-collage"]
