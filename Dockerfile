FROM python:3.7-alpine3.9

ENV NAMU_PORT=3000
ENV NAMU_LANG=en-US

ADD . /app

WORKDIR /app

RUN pip install -r requirements.txt

CMD [ "python", "./app.py" ]