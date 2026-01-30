FROM python:3.12-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1 \
  PYTHONUNBUFFERED=1

RUN apt-get update \
  && apt-get install -y --no-install-recommends gcc \
  && rm -rf /var/lib/apt/lists/*

COPY services/analytics-server-python/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

COPY proto /app/proto
COPY services/analytics-server-python /app/services/analytics-server-python

RUN python /app/services/analytics-server-python/scripts/generate_proto.py

WORKDIR /app/services/analytics-server-python

CMD ["python", "-m", "analytics_server.server"]
