FROM ubuntu:16.04

MAINTAINER Hoto Cocoa <cocoa@hoto.us>

ENV NAMU_PORT=3000
ENV NAMU_LANG=en-US

ADD . /app

WORKDIR /app

RUN apt update && apt install -y --no-install-recommends python3 python3-dev python3-pip python3-setuptools
RUN python3 -m pip install pip --upgrade && python3 -m pip install -r requirements.txt

CMD python3 app.py
