FROM --platform=amd64 golang:1.22

RUN go install github.com/cespare/reflex@latest

ENTRYPOINT ["bash", "watch.sh"]