package config

import "context"

func abortValidateContext() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
