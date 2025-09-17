package codex

import "context"

// Login は OpenAI の API Key を使用して Codex の認証を行います。
func (codex *Codex) Login(ctx context.Context, apiKey string) error {
	codex.authMutex.Lock()
	defer codex.authMutex.Unlock()
	cmd, err := codex.command(ctx, "login", "--api-key", apiKey)
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
