package codex

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type invokeOptions struct {
	ApprovalPolicy   *string // untrusted, on-failure, never
	BaseInstructions *string
	Config           map[string]any
	Cwd              *string // 実行時のカレントディレクトリ
	IncludePlanTool  *bool
	Model            *string // モデル
	Profile          *string
	Sandbox          *string // read-only, workspace-write, danger-full-access
}

// InvokeOption は Invoke 呼び出しに渡すオプション関数の型です。
// 各オプションは内部の `invokeOptions` を変更し、エラーを返すことができます。
type InvokeOption func(*invokeOptions) error

// WithApprovalPolicy は承認ポリシー（"untrusted", "on-failure", "never"）
// を設定するオプションを返します。無効な値が与えられるとエラーになります。
func WithApprovalPolicy(approvalPolicy string) InvokeOption {
	return func(o *invokeOptions) error {
		switch approvalPolicy {
		default:
			return errors.New("invalid approval-policy: " + approvalPolicy)
		case "untrusted", "on-failure", "never":
			o.ApprovalPolicy = &approvalPolicy
			return nil
		}
	}
}

// WithBaseInstructions はツールに渡すベース命令（base instructions）を設定する
// オプションを返します。プロンプトに常に付加したい命令を指定します。
func WithBaseInstructions(baseInstructions string) InvokeOption {
	return func(o *invokeOptions) error {
		o.BaseInstructions = &baseInstructions
		return nil
	}
}

// WithConfig はツール呼び出し時に渡す追加の設定マップを指定するオプションを返します。
func WithConfig(config map[string]any) InvokeOption {
	return func(o *invokeOptions) error {
		o.Config = config
		return nil
	}
}

// WithCwd はツール実行時のカレントディレクトリを指定するオプションを返します。
func WithCwd(cwd string) InvokeOption {
	return func(o *invokeOptions) error {
		o.Cwd = &cwd
		return nil
	}
}

// WithIncludePlanTool は実行時にプランツールを含めるかどうかを設定するオプションを返します。
func WithIncludePlanTool(includePlanTool bool) InvokeOption {
	return func(o *invokeOptions) error {
		o.IncludePlanTool = &includePlanTool
		return nil
	}
}

// WithModel は使用するモデル名を指定するオプションを返します。
func WithModel(model string) InvokeOption {
	return func(o *invokeOptions) error {
		o.Model = &model
		return nil
	}
}

// WithProfile は実行時に使用するプロファイル名を設定するオプションを返します。
func WithProfile(profile string) InvokeOption {
	return func(o *invokeOptions) error {
		o.Profile = &profile
		return nil
	}
}

// WithSandbox はサンドボックス設定（"read-only", "workspace-write", "danger-full-access"）
// を指定するオプションを返します。無効な値が与えられるとエラーになります。
func WithSandbox(sandbox string) InvokeOption {
	return func(o *invokeOptions) error {
		switch sandbox {
		default:
			return errors.New("invalid sandbox: " + sandbox)
		case "read-only", "workspace-write", "danger-full-access":
			o.Sandbox = &sandbox
			return nil
		}
	}
}

// Invoke は Codex を実行して結果を返します。
// 指定可能なオプションの詳細は以下を参照してください。
// https://github.com/openai/codex/blob/main/docs/advanced.md#codex-mcp-server-quickstart
func (codex *Codex) Invoke(ctx context.Context, prompt string, options ...InvokeOption) (string, error) {
	opts := invokeOptions{}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return "", err
		}
	}

	cmd, err := codex.command(ctx, "mcp")
	if err != nil {
		return "", err
	}

	if opts.Cwd != nil {
		cmd.Dir = *opts.Cwd
	}

	client := mcp.NewClient(&mcp.Implementation{}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		return "", err
	}
	defer session.Close()

	arguments := map[string]any{
		"prompt": prompt,
	}

	if opts.ApprovalPolicy != nil {
		arguments["approval-policy"] = *opts.ApprovalPolicy
	}

	if opts.BaseInstructions != nil {
		arguments["base-instructions"] = *opts.BaseInstructions
	}

	if opts.Config != nil {
		arguments["config"] = opts.Config
	}

	if opts.IncludePlanTool != nil {
		arguments["include-plan-tool"] = *opts.IncludePlanTool
	}

	if opts.Model != nil {
		arguments["model"] = *opts.Model
	}

	if opts.Profile != nil {
		arguments["profile"] = *opts.Profile
	}

	if opts.Sandbox != nil {
		arguments["sandbox"] = *opts.Sandbox
	}

	params := &mcp.CallToolParams{
		Name:      "codex",
		Arguments: arguments,
	}
	res, err := session.CallTool(ctx, params)
	if err != nil {
		return "", err
	}

	text := res.Content[0].(*mcp.TextContent).Text
	if res.IsError {
		return "", errors.New(text)
	}

	return text, nil
}
