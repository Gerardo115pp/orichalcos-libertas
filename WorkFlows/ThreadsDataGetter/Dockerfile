FROM python:rc-alpine3.13

WORKDIR /thread-extractor

RUN apk add gcc musl-dev libffi-dev openssl-dev python3-dev
RUN pip install bs4 prefect

COPY . .

CMD ["python", "threads_extractor.py"]