FROM python:3.10-alpine

WORKDIR /app

COPY bot bot
COPY common common
RUN mv bot/bot.py bot.py

RUN apk add build-base

RUN pip3 install -r bot/requirements.txt

CMD ["python3", "bot.py"]
