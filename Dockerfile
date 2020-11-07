############################
# STEP 1 build executable binary
############################
FROM golang:alpine
RUN apk update && apk add --no-cache git && apk add bash
WORKDIR $GOPATH/trackerapp/

COPY . .

# Build the binary.
#RUN cd ./cmd/tracker/
#RUN go get -d -v
RUN go build ./cmd/tracker/main.go

ENTRYPOINT ["go", "run", "./cmd/tracker/main.go"]