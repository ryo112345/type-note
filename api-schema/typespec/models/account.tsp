import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

namespace TypeNote.Models;

// アカウント情報
model Account {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  fullName: string;
  thumbnail?: string;
  lastLoginAt: utcDateTime;
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

// OAuth認証リクエスト
model CreateOrGetAccountRequest {
  email: string;
  firstName: string;
  lastName: string;
  provider: string;
  providerAccountId: string;
  thumbnail?: string;
}

// アカウントレスポンス
model AccountResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  fullName: string;
  thumbnail?: string;
  lastLoginAt: utcDateTime;
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

// 簡易アカウント情報
model AccountSummary {
  id: string;
  firstName: string;
  lastName: string;
  thumbnail?: string;
}
