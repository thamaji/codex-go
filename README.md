# codex-go

## 概要

`codex-go` は OpenAI が提供する CLI 型のコーディングエージェント `codex` を
Go から簡単に利用するためのラッパーライブラリです。

本ライブラリは `codex` v0.36.0 での動作を確認しています。

## インストール

```bash
go get github.com/thamaji/codex-go
```

## 使用例

簡単な使用例です。環境変数 `OPENAI_APIKEY` に API キーを入れて実行します。

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/thamaji/codex-go"
)

func main() {
    apiKey := os.Getenv("OPENAI_APIKEY")

    client := codex.New(codex.WithLogger(os.Stderr, "info"))

    ctx := context.Background()
    if err := client.Login(ctx, apiKey); err != nil {
        log.Fatal(err)
    }

    // シンプルな呼び出し例
    text, err := client.Invoke(
        ctx,
        "hello",
        codex.WithCwd("."),
        codex.WithSandbox("read-only"),
        codex.WithApprovalPolicy("never"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(text)
}
```
