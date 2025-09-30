# Google OAuth 2.0 Setup Guide

このガイドでは、Todo AppでGoogle OAuth認証を使用するために必要なGoogle Cloud Consoleの設定手順を説明します。

## 📋 必要な設定項目

### 1. Google Cloud Console でのプロジェクト作成

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. 新規プロジェクトを作成または既存プロジェクトを選択
3. プロジェクト名: `todo-app-oauth` (任意)

### 2. OAuth同意画面の設定

**パス**: `APIs & Services > OAuth consent screen`

```
User Type: External
App name: Todo App
User support email: <your-email@gmail.com>
Developer contact information: <your-email@gmail.com>
Authorized domains: localhost (開発時)
```

### 3. OAuth 2.0認証情報の作成

**パス**: `APIs & Services > Credentials > Create Credentials > OAuth 2.0 Client IDs`

```
Application type: Web application
Name: Todo App OAuth Client

Authorized JavaScript origins:
- http://localhost:3000

Authorized redirect URIs:
- http://localhost:8080/api/v1/auth/google/callback
```

### 4. 認証情報の設定

作成後に表示される認証情報を `.env` ファイルに設定：

```bash
# .envファイルを編集
GOOGLE_CLIENT_ID=<取得したClient ID>
GOOGLE_CLIENT_SECRET=<取得したClient Secret>
```

## 🔧 環境変数の例

```env
# Google OAuth 2.0 Configuration
GOOGLE_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-abcdefghijklmnopqrstuvwxyz
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/auth/google/callback

# JWT Configuration
JWT_SECRET=<32文字以上の強力なランダム文字列>
JWT_EXPIRES_HOURS=24
```

## 🚨 セキュリティ注意事項

### 開発環境
- `.env`ファイルは `.gitignore` に追加済み
- `localhost` での開発は安全

### 本番環境
- 強力なJWTシークレットキーを使用
- HTTPS必須 (`SESSION_COOKIE_SECURE=true`)
- 本番ドメインを認証済みドメインに追加
- 環境変数は安全な方法で管理

## 🔍 動作確認

1. バックエンドサーバー起動: `go run main.go`
2. フロントエンドサーバー起動: `npm start`
3. `http://localhost:3000/login` にアクセス
4. "Sign in with Google" ボタンをテスト

## ❌ よくあるエラー

### `redirect_uri_mismatch`
**原因**: Redirect URIが一致しない
**解決**: Google Cloud Consoleで正確なURIを設定

### `invalid_client`
**原因**: Client IDまたはSecretが間違い
**解決**: `.env`ファイルの値を確認

### `access_blocked`
**原因**: OAuth同意画面の設定不備
**解決**: User TypeをExternalに設定し、必要項目を入力

## 📚 参考リンク

- [Google OAuth 2.0 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth 2.0 Scopes](https://developers.google.com/identity/protocols/oauth2/scopes)