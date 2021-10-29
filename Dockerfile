FROM alpine

RUN apk add --no-cache go=1.17.1
RUN cp -r . .

ENTRYPOINT ["go", "run", "main.go"]