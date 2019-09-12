FROM alpine
RUN apk add --no-cache ca-certificates
COPY . /
EXPOSE 80
ENTRYPOINT ["./summary"]
