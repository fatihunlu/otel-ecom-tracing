var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var port = Environment.GetEnvironmentVariable("HTTP_PORT") ?? "5001";
app.Urls.Add($"http://0.0.0.0:{port}");

app.MapGet("/health", () => new { status = "ok", service = "product-api" });

app.Run();