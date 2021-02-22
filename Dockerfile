FROM golang:1.15.7 AS base

ENV PORT 8000

COPY . /app
WORKDIR /app

RUN make server

FROM python:3.8-buster AS dist

COPY --from=base /app/app.out /
COPY --from=base /app/web /web

RUN pip install youtube_dl

CMD PORT=${PORT} ./app.out