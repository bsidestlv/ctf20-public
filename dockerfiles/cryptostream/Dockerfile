FROM golang:1.14 as builder
ENV FLAG BSidesTLV{refractions of dawn in the andromeda universe}
ENV PORT 8080
#ENV HTTPS ""

RUN mkdir /app
ADD . /app
WORKDIR /app
EXPOSE $PORT

ENV GOPATH=/app
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o cryptostream main.go

FROM scratch
COPY --from=builder /app/cryptostream ./
COPY --from=builder /app/server.* ./

CMD ["./cryptostream"]
