FROM python:3.10
ENV TARGET=localhost:8080
WORKDIR /usr/src/app

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD [ "python", "clients.py" ]