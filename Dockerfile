FROM python:3.10.5-alpine

MAINTAINER 2du <min08101@naver.com>
MAINTAINER hoparkgo9ma <me@ho9.me>

ENV NAMU_DB_TYPE sqlite
ENV NAMU_DB data
ENV NAMU_HOST 0.0.0.0
ENV NAMU_PORT 3000
ENV NAMU_LANG en-US
ENV NAMU_MARKUP namumark
ENV NAMU_ENCRYPT sha3
ENV NAMU_DOCKER O

ADD . /app
WORKDIR /app

RUN pip install -r requirements.txt
EXPOSE 3000

CMD [ "python", "./app.py" ]