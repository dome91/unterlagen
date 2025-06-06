FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY unterlagen .
COPY platform/web/assets ./platform/web/assets

EXPOSE 8080

CMD ["./unterlagen"]