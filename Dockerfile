FROM python:3

ENV NAMU_DB=data
ENV NAMU_HOST=0.0.0.0
ENV NAMU_PORT=3000
ENV NAMU_LANG=en-US

ADD . /app

WORKDIR /app

RUN pip install -r requirements.txt

CMD [ "python", "./app.py" ]
