import os
from opentelemetry import trace
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor

def init_tracer():
    svc_name    = os.getenv("SERVICE_NAME", "ecom.payment.api")
    svc_version = os.getenv("SERVICE_VERSION", "0.1.0")
    env         = os.getenv("ENVIRONMENT", "local")

    resource = Resource.create({
        "service.name": svc_name,
        "service.version": svc_version,
        "deployment.environment": env,
    })

    endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector:4317")
    endpoint = endpoint.replace("http://", "").replace("https://", "")

    exporter = OTLPSpanExporter(endpoint=endpoint, insecure=True)

    provider = TracerProvider(resource=resource)
    provider.add_span_processor(BatchSpanProcessor(exporter))
    trace.set_tracer_provider(provider)
