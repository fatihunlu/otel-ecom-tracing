using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Metrics;

var builder = WebApplication.CreateBuilder(args);

// --- OpenTelemetry config ---
var svcName = Environment.GetEnvironmentVariable("SERVICE_NAME") ?? "ecom.product.api";
var svcVersion = Environment.GetEnvironmentVariable("SERVICE_VERSION") ?? "0.1.0";
var environment = Environment.GetEnvironmentVariable("ENVIRONMENT") ?? "local";

builder.Services.AddOpenTelemetry()
    .ConfigureResource(r => r
        .AddService(serviceName: svcName, serviceVersion: svcVersion)
        .AddAttributes(new KeyValuePair<string, object>[]
        {
            new("deployment.environment", environment)
        }))
    .WithTracing(t => t
        .AddAspNetCoreInstrumentation(opt =>
        {
            opt.RecordException = true;
            opt.EnrichWithHttpRequest = (activity, request) =>
            {
                activity.SetTag("http.request_content_length", request.ContentLength);
            };
            opt.EnrichWithHttpResponse = (activity, response) =>
            {
                activity.SetTag("http.response_content_length", response.ContentLength);
            };
        })
        .AddHttpClientInstrumentation(opt => opt.RecordException = true)
        .AddOtlpExporter()
    )
    .WithMetrics(m => m
        .AddRuntimeInstrumentation()
        .AddMeter("Microsoft.AspNetCore")
        .AddOtlpExporter()
    );

var app = builder.Build();

var port = Environment.GetEnvironmentVariable("HTTP_PORT") ?? "5001";
app.Urls.Add($"http://0.0.0.0:{port}");

app.MapGet("/health", () => new { status = "ok", service = "product-api" });
app.MapGet("/products", () => new[] { "Keyboard", "Mouse", "Monitor" });

app.Run();