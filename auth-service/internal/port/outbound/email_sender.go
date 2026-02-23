package outbound

import "context"

type EmailSender interface {
	SendPasswordResetEmail(ctx context.Context, toEmail, resetLink string) error
}
