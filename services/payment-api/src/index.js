const express = require("express");
const app = express();

const PORT = process.env.HTTP_PORT || 5003;

app.get("/health", (_req, res) => {
  res.json({ status: "ok", service: "payment-api" });
});

app.listen(PORT, () => {
  console.log(`payment-api listening on :${PORT}`);
});
