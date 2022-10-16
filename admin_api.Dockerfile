FROM python:3.10

WORKDIR /app

COPY admin_api admin_api
COPY common common
RUN mv admin_api/main.py main.py

RUN pip3 install -r admin_api/requirements.txt

CMD [ "uvicorn", "main:app", "--host", "0.0.0.0"]
