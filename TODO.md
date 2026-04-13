# post2bsky — TODO

## 概要
ターミナルから Bluesky に投稿だけできる TUI アプリ。
他の人の投稿は見れない。自分の直近投稿（最大5件）のみ参考表示。

---

## 使用技術
- 言語: Go
- TUI: Bubble Tea + Bubbles + Lipgloss
- 認証: App Password（環境変数 or 設定ファイル）+ セッショントークンのキャッシュ
- Bluesky API: 直接 HTTP（`net/http`）

## 依存パッケージ
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/bubbles`
- `github.com/charmbracelet/lipgloss`
- `github.com/BurntSushi/toml`

---

## ファイル構成

```
zenn-post2bsky/
├── main.go
├── go.mod
├── config/
│   └── config.go
├── bsky/
│   ├── client.go
│   ├── session.go      # CreateSession(), RefreshSession(), SaveSession(), LoadSession()
│   ├── post.go
│   └── feed.go
└── ui/
    ├── model.go
    ├── vim.go          # Vim モード管理（Normal/Insert）
    ├── keys.go
    ├── styles.go
    └── messages.go
```

---

## タスク

- [ ] `go mod init` + 依存パッケージ取得
- [ ] `config/config.go` — 環境変数 / TOML ファイルから認証情報を読み込む
- [ ] `bsky/client.go` — `doXRPC()` HTTP ヘルパー、`Client` struct
- [ ] `bsky/session.go` — セッション管理（認証・更新・キャッシュ保存/読み込み）
- [ ] `bsky/post.go` — `CreatePost(text string)`
- [ ] `bsky/feed.go` — `GetOwnRecentPosts(limit int)`
- [ ] `ui/messages.go` — カスタム `tea.Msg` 型定義
- [ ] `ui/vim.go` — Vim モード（Normal/Insert）と Normal モードのキー処理
- [ ] `ui/keys.go` — キーバインド定義
- [ ] `ui/styles.go` — Lipgloss スタイル定義
- [ ] `ui/model.go` — `Model`, `Init`, `Update`, `View`, コマンドファクトリ
- [ ] `main.go` — エントリポイント

---

## セッション永続化

起動のたびにログインしない。トークンをキャッシュし、次回起動時に再利用する。

### キャッシュファイル
`~/.config/post2bsky/session.json`

```json
{
  "access_jwt": "...",
  "refresh_jwt": "...",
  "did": "did:plc:...",
  "handle": "yourhandle.bsky.social"
}
```

### 起動時の認証フロー
```
キャッシュあり
  └─ refreshSession 試行
       ├─ 成功 → 新トークンをキャッシュに保存して起動
       └─ 失敗 → createSession（App Password）→ キャッシュに保存して起動

キャッシュなし
  └─ createSession（App Password）→ キャッシュに保存して起動
```

`accessJwt` の有効期限は約2時間。`refreshJwt` は約90日。
通常の使い方であれば `refreshSession` だけでほぼ済む。

---

## Vim ライクなテキスト編集

### モード

| モード | 説明 |
|--------|------|
| Normal | カーソル移動・操作コマンド入力（起動時はこちら） |
| Insert | テキスト入力 |

### Normal モードのキーバインド

| キー | 動作 |
|------|------|
| `i` | カーソル位置の前に挿入（Insert モードへ） |
| `a` | カーソル位置の後に挿入（Insert モードへ） |
| `o` | 次の行に新規行を追加して挿入（Insert モードへ） |
| `O` | 前の行に新規行を追加して挿入（Insert モードへ） |
| `h` / `l` | 左 / 右 |
| `j` / `k` | 下 / 上 |
| `0` | 行頭へ |
| `$` | 行末へ |
| `w` | 次の単語へ |
| `b` | 前の単語へ |
| `dd` | 行削除 |
| `u` | アンドゥ |
| `Ctrl+S` | 投稿（Normal / Insert どちらでも可） |
| `Ctrl+C` | 終了 |

### Insert モードのキーバインド

| キー | 動作 |
|------|------|
| `Esc` | Normal モードへ戻る |
| `Ctrl+S` | 投稿 |
| 通常の文字入力 | テキスト入力 |

### 実装方針
- `bubbles/textarea` の内部テキストバッファはそのまま利用する
- `ui/vim.go` に `VimMode` 型と Normal モードのキー処理関数を定義する
- `ui/model.go` の `Update()` で現在のモードに応じてキー処理を分岐する
- モードインジケータ（`-- INSERT --` / `-- NORMAL --`）をフッターに表示する

---

## 設定方法

### 環境変数（優先）
```bash
export BSKY_HANDLE=yourhandle.bsky.social
export BSKY_APP_PASSWORD=xxxx-xxxx-xxxx-xxxx
```

### 設定ファイル（フォールバック）
```toml
# ~/.config/post2bsky/config.toml
handle = "yourhandle.bsky.social"
app_password = "xxxx-xxxx-xxxx-xxxx"
# pds = "https://bsky.social"  # optional
```

---

## TUI レイアウト

```
╭─────────────────────────────────────╮
│  post2bsky                          │
│  ─────────────────────────────────  │
│  [textarea: 5行, 300文字制限]        │
│                                     │
│                               0/300 │
│  ─────────────────────────────────  │
│  -- NORMAL --        ctrl+s: Post   │
│  ─────────────────────────────────  │
│  Recent posts (last 5):             │
│  • 投稿テキスト...                   │
│  • 別の投稿...                       │
╰─────────────────────────────────────╯
```

---

## Bluesky API Endpoints

| 操作 | エンドポイント | ホスト |
|------|--------------|-------|
| 認証 | `POST /xrpc/com.atproto.server.createSession` | PDS (`bsky.social`) |
| トークン更新 | `POST /xrpc/com.atproto.server.refreshSession` | 同上 |
| 投稿 | `POST /xrpc/com.atproto.repo.createRecord` | 同上 |
| 自分の投稿取得 | `GET /xrpc/app.bsky.feed.getAuthorFeed?actor=<DID>&limit=5&filter=posts_no_replies` | `api.bsky.app` |

---

## 動作確認

```bash
# 起動
BSKY_HANDLE=yourhandle.bsky.social BSKY_APP_PASSWORD=xxxx go run .

# ビルド
go build -o post2bsky .
./post2bsky
```

### 確認項目
- [ ] 初回起動時に認証成功し、セッションがキャッシュに保存される
- [ ] 2回目以降の起動でキャッシュトークンが使われ、即座に起動する
- [ ] `i` キーで Insert モードに入り、テキストが入力できる
- [ ] `Esc` で Normal モードに戻り、`hjkl` でカーソル移動できる
- [ ] `dd` で行削除、`u` でアンドゥができる
- [ ] 文字カウンターが入力に応じて更新される
- [ ] 300文字超えると `Ctrl+S` で投稿できない
- [ ] `Ctrl+S` で投稿でき、成功後にテキストがクリアされ投稿リストが更新される
- [ ] `Ctrl+C` でアプリが終了する
