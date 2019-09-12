FROM alpine
RUN apk add --no-cache ca-certificates
COPY . /
EXPOSE 443
 ENTRYPOINT [ "./gateway" ]
