import time
from fastapi import FastAPI
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry import trace
from .otel import init_tracer

init_tracer()

app = FastAPI(title="payment-api")
FastAPIInstrumentor().instrument_app(app)

@app.get("/health")
def health():
    return {"status": "ok", "service": "payment-api"}

@app.get("/pay")
def pay():
    with trace.get_tracer("payment").start_as_current_span("charge"):
        time.sleep(0.05)
    return {"ok": True}
