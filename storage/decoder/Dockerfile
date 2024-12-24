FROM --platform=linux/arm64 arm64v8/python:3.11.9-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY ./requirements.txt ./requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

RUN rm -rf ./requirements.txt

COPY ./domain ./domain

COPY ./preload.py ./

RUN python3 preload.py

RUN rm -rf ./preload.py

COPY ./main.py ./

COPY ./adapter ./adapter

COPY ./service ./service

CMD ["python3", "main.py"]
