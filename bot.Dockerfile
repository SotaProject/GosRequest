FROM python:3.10-alpine

WORKDIR /app

COPY bot bot
COPY common common
RUN mv bot/main.py main.py

RUN pip3 install -r bot/requirements.txt

CMD ["python3", "main.py"]
