import "@typespec/http";
import "@typespec/openapi3";
import "./models/common.tsp";
import "./models/account.tsp";
import "./routes/accounts.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service({
  title: "Type Note API",
})
@server("http://localhost:8080", "Development server")
namespace TypeNote;
