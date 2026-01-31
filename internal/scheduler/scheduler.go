package scheduler

import (
	"fmt"
	"strings"
	"time"

	"account-manager/internal/models"
	"account-manager/internal/repository"
	"account-manager/internal/service"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron         *cron.Cron
	accountRepo  *repository.AccountRepository
	emailRepo    *repository.EmailRepository
	emailService *service.EmailService
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:         cron.New(),
		accountRepo:  repository.NewAccountRepository(),
		emailRepo:    repository.NewEmailRepository(),
		emailService: service.NewEmailService(),
	}
}

func (s *Scheduler) Start() {
	// Run expiry check every hour
	s.cron.AddFunc("0 * * * *", s.CheckExpiringAccounts)

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) CheckExpiringAccounts() {
	sysConfig, err := s.emailRepo.GetSystemConfig()
	if err != nil {
		return
	}

	daysBefore := sysConfig.ReminderDaysBefore
	accounts, err := s.accountRepo.FindExpiringAccounts(daysBefore)
	if err != nil || len(accounts) == 0 {
		return
	}

	// Send reminder email
	err = s.sendExpiryReminder(accounts, daysBefore)
	if err != nil {
		return
	}

	// Mark as reminder sent
	var ids []uint
	for _, acc := range accounts {
		ids = append(ids, acc.ID)
	}
	s.accountRepo.MarkReminderSent(ids)
}

func (s *Scheduler) sendExpiryReminder(accounts []models.Account, daysBefore int) error {
	subject := fmt.Sprintf("账号过期提醒 - %d个账号即将在%d天后过期", len(accounts), daysBefore)

	var tableRows strings.Builder
	for _, acc := range accounts {
		expireDate := ""
		if acc.ExpireAt != nil {
			expireDate = acc.ExpireAt.Format("2006-01-02")
		}
		tableRows.WriteString(fmt.Sprintf(`
			<tr>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
				<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
			</tr>
		`, acc.Account, acc.AccountType, expireDate))
	}

	content := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2 style="color: #1890ff;">账号过期提醒</h2>
		<p>以下 <strong>%d</strong> 个账号将在 <strong>%d</strong> 天后过期，请及时处理：</p>
		<table style="border-collapse: collapse; width: 100%%;">
			<thead>
				<tr style="background-color: #f5f5f5;">
					<th style="padding: 8px; border: 1px solid #ddd; text-align: left;">账号</th>
					<th style="padding: 8px; border: 1px solid #ddd; text-align: left;">类型</th>
					<th style="padding: 8px; border: 1px solid #ddd; text-align: left;">过期日期</th>
				</tr>
			</thead>
			<tbody>
				%s
			</tbody>
		</table>
		<p style="color: #666; margin-top: 20px;">
			发送时间: %s<br>
			此邮件由账号管理系统自动发送
		</p>
	</body>
	</html>
	`, len(accounts), daysBefore, tableRows.String(), time.Now().Format("2006-01-02 15:04:05"))

	return s.emailService.SendEmail(subject, content)
}

// ManualCheck allows manual triggering of expiry check
func (s *Scheduler) ManualCheck() (int, error) {
	sysConfig, err := s.emailRepo.GetSystemConfig()
	if err != nil {
		return 0, err
	}

	daysBefore := sysConfig.ReminderDaysBefore
	accounts, err := s.accountRepo.FindExpiringAccounts(daysBefore)
	if err != nil {
		return 0, err
	}

	if len(accounts) == 0 {
		return 0, nil
	}

	err = s.sendExpiryReminder(accounts, daysBefore)
	if err != nil {
		return 0, err
	}

	var ids []uint
	for _, acc := range accounts {
		ids = append(ids, acc.ID)
	}
	s.accountRepo.MarkReminderSent(ids)

	return len(accounts), nil
}
