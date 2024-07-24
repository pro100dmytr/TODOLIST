FROM ubuntu:latest
LABEL authors="dimon"

ENTRYPOINT ["top", "-b"]