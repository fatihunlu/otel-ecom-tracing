# payment-api (Python · FastAPI · OpenTelemetry)

Basit bir ödeme servisi. FastAPI ile HTTP endpoint’leri sunar, OpenTelemetry ile **OTLP/gRPC** üzerinden **OpenTelemetry Collector**’a trace gönderir. Jaeger ile izleme, Grafana/Prometheus ile metrik takibi yapılır.

## Çalıştırma
```bash
docker compose -f infra/docker-compose.yml --env-file .env build payment-api
docker compose -f infra/docker-compose.yml --env-file .env up -d payment-api
