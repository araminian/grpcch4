FROM golang:1.23.3-alpine AS base
WORKDIR /app
COPY . .
RUN cd server && go build -o server .

FROM alpine:latest
WORKDIR /app
USER 1001
COPY --from=base /app/server .
EXPOSE 5557
ENTRYPOINT ["./server"]
