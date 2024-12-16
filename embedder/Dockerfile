FROM --platform=linux/arm64 arm64v8/python:3.13.1-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY ./requirements.txt ./requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

RUN rm -rf ./requirements.txt

COPY ./main.py ./main.py

CMD ["python3", "main.py"]
