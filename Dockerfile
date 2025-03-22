FROM ubuntu:jammy

RUN mkdir -pv /data/logs/
WORKDIR /app
ENTRYPOINT [ "/app/main" ]