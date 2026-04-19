# post2bsky

ターミナルから Bluesky に投稿するだけの TUI アプリ。

```
╭─────────────────────────────────────╮
│  post2bsky                          │
│  ─────────────────────────────────  │
│  What's on your mind?               │
│                                     │
│                               0/300 │
│  ─────────────────────────────────  │
│                    ctrl+s: Post   ctrl+c: Quit │
╰─────────────────────────────────────╯
```

## Getting Started

### 1. インストール

**Homebrew（macOS / Linux）**

```bash
brew install lCyou/tap/post2bsky
```

> `brew tap lCyou/tap && brew install post2bsky` でも同じです。

**直接ダウンロード**

[Releases](https://github.com/lCyou/zenn-post2bsky/releases/latest) から環境に合ったバイナリをダウンロードして解凍し、PATH の通った場所に置いてください。

| ファイル | 環境 |
|---|---|
| `post2bsky_darwin_arm64.tar.gz` | macOS (Apple Silicon) |
| `post2bsky_darwin_amd64.tar.gz` | macOS (Intel) |
| `post2bsky_linux_amd64.tar.gz` | Linux (x86_64) |
| `post2bsky_linux_arm64.tar.gz` | Linux (ARM) |
| `post2bsky_windows_amd64.zip` | Windows |

**Go がある場合**

```bash
go install github.com/kyou/post2bsky@latest
```

### 2. App Password を発行する

Bluesky の通常パスワードは使用しません。専用の App Password を発行してください。

1. [bsky.app](https://bsky.app) にログイン
2. 設定 → **App Passwords** → **Add App Password**
3. 名前を入力（例: `post2bsky`）して作成
4. 表示された `xxxx-xxxx-xxxx-xxxx` 形式のパスワードをコピー

### 3. 起動する

```bash
post2bsky
```

初回起動時にハンドルと App Password を聞かれます。入力すると次回から自動でログインします。

```
Bluesky Handle: yourhandle.bsky.social
App Password:   ················
```

## 使い方

| キー | 動作 |
|---|---|
| `Ctrl+S` | 投稿 |
| `Ctrl+C` | 終了 |

300文字以内で入力して `Ctrl+S` で投稿できます。

## 設定

認証情報は `~/.config/post2bsky/config.toml` に保存されます。環境変数でも指定できます。

```bash
export BSKY_HANDLE=yourhandle.bsky.social
export BSKY_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx
```
