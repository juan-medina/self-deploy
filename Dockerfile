FROM alpine:latest

COPY self-deploy self-deploy

CMD ["./self-deploy"]