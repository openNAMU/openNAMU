FROM python:alpine

LABEL org.opencontainers.image.authors="2du <min08101@naver.com>, hoparkgo9ma <me@ho9.me>"

ENV NAMU_DB_TYPE sqlite
ENV NAMU_DB data
ENV NAMU_HOST 0.0.0.0
ENV NAMU_PORT 3000
ENV NAMU_LANG ko-KR
ENV NAMU_MARKUP namumark
ENV NAMU_ENCRYPT sha3

ADD . /app
WORKDIR /app

RUN pip install -r requirements.txt
EXPOSE 3000

CMD [ "python", "./app.py" ]