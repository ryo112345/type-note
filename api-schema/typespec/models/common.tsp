import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

namespace TypeNote.Models;

// 共通エラーレスポンス
@error
model ErrorResponse {
  code: string;
  message: string;
  details?: unknown;
}

// Not Found エラー
@error
model NotFoundError {
  code: "NOT_FOUND";
  message: string;
}

// Unauthorized エラー
@error
model UnauthorizedError {
  code: "UNAUTHORIZED";
  message: string;
}

// Forbidden エラー
@error
model ForbiddenError {
  code: "FORBIDDEN";
  message: string;
}

// Bad Request エラー
@error
model BadRequestError {
  code: "BAD_REQUEST";
  message: string;
  details?: unknown;
}

// 成功レスポンス
model SuccessResponse {
  success: boolean;
}
