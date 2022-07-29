package mail

import (
	"context"
	"fmt"

	"github.com/ITA-Dnipro/Dp-230-Report-Service/internal/model"
)

type MailSender interface {
	Send(ctx context.Context, mail model.Mail) error
}

type ReportMailSender struct {
	mailSender     MailSender
	serviceBaseURL string
}

func NewReportMailSender(ms MailSender, url string) *ReportMailSender {
	return &ReportMailSender{mailSender: ms, serviceBaseURL: url}
}

func (r *ReportMailSender) Send(ctx context.Context, report model.Report) error {
	body := fmt.Sprintf("<p>Test results <b>ready</b> at <a href='%s/report/%s'>report</a>.</p>", r.serviceBaseURL, report.ID)
	mail := model.Mail{
		To:      []string{report.Email},
		Subject: "scan report",
		Body:    body,
	}
	return r.mailSender.Send(ctx, mail)
}
