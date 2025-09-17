package codex

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// CodexOption は Codex の設定を変更するためのオプション関数です。
// これを使って `New` 呼び出し時に実装固有の設定を注入できます。
type CodexOption func(*Codex)

// WithExecutablePath は Codex が使用する実行ファイルのパスを設定する
// オプションを返します。デフォルトは `codex` です。
func WithExecutablePath(path string) CodexOption {
	return func(codex *Codex) {
		codex.SetExecutablePath(path)
	}
}

// WithLogger は Codex のログ出力先とログレベルを設定するオプションを返します。
// `w` にログが書き出され、`level` でログの詳細度を制御します。
func WithLogger(w io.Writer, level string) CodexOption {
	return func(codex *Codex) {
		codex.SetLogger(w, level)
	}
}

// Codex は codex コマンドをラップするクライアント構造体です。
// 内部で実行コマンドのパスやログ設定、認証用のロックを保持します。
type Codex struct {
	authMutex sync.Mutex

	executablePath string    // 実行コマンドのパス（デフォルト：codex）
	logWriter      io.Writer // ログの出力先（デフォルト：nil）
	logLevel       string    // ログレベル（デフォルト：info、有効な値：error, warn, info, debug, trace, off）
}

// New は Codex のインスタンスを作成します。
// WithExecutablePath で codex コマンドの実行パスを指定します。
// WithLogger で codex コマンドのログ出力を設定します。
func New(options ...CodexOption) *Codex {
	codex := Codex{
		executablePath: "codex",
		logWriter:      nil,
		logLevel:       "info",
	}
	for _, opt := range options {
		opt(&codex)
	}
	return &codex
}

// SetExecutablePath は Codex インスタンスの実行ファイルパスを設定します。
// テストやカスタムビルドを使う場合に利用します。
func (codex *Codex) SetExecutablePath(path string) {
	codex.executablePath = path
}

// SetLogger は Codex インスタンスのログ出力先とログレベルを設定します。
// `WithLogger` と同様の振る舞いを持ち、インスタンス生成後に設定を変更できます。
func (codex *Codex) SetLogger(w io.Writer, level string) {
	codex.logWriter = w
	codex.logLevel = level
}

func (codex *Codex) command(ctx context.Context, arg ...string) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, codex.executablePath, arg...)
	if codex.logWriter != nil {
		switch codex.logLevel {
		default:
			return nil, errors.New("invalid log level")
		case "error", "warn", "info", "debug", "trace", "off":
			cmd.Stderr = codex.logWriter
			cmd.Env = append([]string{fmt.Sprintf("RUST_LOG=codex_core=%s,codex_tui=%s", codex.logLevel, codex.logLevel)}, os.Environ()...)
		}
	}
	return cmd, nil
}
