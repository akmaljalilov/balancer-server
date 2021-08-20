FROM golang:1.15.3 as builder
COPY . /go/src/velox/server
WORKDIR /go/src/velox/server
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
ENV buildNumber=BUILD_NUMBER
ENV branchName=BRANCH
RUN apk --no-cache add ca-certificates bash
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/velox/server .

CMD ["./server"]
