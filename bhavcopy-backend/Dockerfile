FROM golang:1.14.9-alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/bhavcopy-backend /app/
COPY Data/ /app/Data
USER root
RUN chmod 777 /app/Data
COPY config/ /app/config
WORKDIR /app
CMD ["./bhavcopy-backend"]

