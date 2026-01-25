package handlers

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// SMTPVerifier 用于验证邮箱的 SMTP 验证器
type SMTPVerifier struct {
	fromEmail string
	timeout   time.Duration
}

// NewSMTPVerifier 创建新的 SMTP 验证器
func NewSMTPVerifier() *SMTPVerifier {
	return &SMTPVerifier{
		fromEmail: "verify@example.com", // 用于 MAIL FROM 命令
		timeout:   10 * time.Second,
	}
}

// VerifyEmail 验证单个邮箱地址
// 返回状态: "live", "dead", "unknown"
func (v *SMTPVerifier) VerifyEmail(email string) (string, error) {
	// 1. 验证邮箱格式
	if !strings.Contains(email, "@") {
		return "dead", fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "dead", fmt.Errorf("invalid email format")
	}
	domain := parts[1]

	// 2. 查询 MX 记录
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return "dead", fmt.Errorf("no MX records found for domain: %s", domain)
	}

	// 3. 尝试多个 MX 服务器和端口
	var lastErr error
	ports := []string{"25", "587", "465"} // SMTP, Submission, SMTPS

	for _, mxRecord := range mxRecords {
		mxHost := strings.TrimSuffix(mxRecord.Host, ".")

		for _, port := range ports {
			status, err := v.tryVerifyWithPort(email, mxHost, port)
			if err == nil {
				return status, nil
			}
			lastErr = err
		}
	}

	// 所有尝试都失败了
	return "unknown", fmt.Errorf("cannot verify email: %v", lastErr)
}

// tryVerifyWithPort 尝试使用指定端口验证邮箱
func (v *SMTPVerifier) tryVerifyWithPort(email, mxHost, port string) (string, error) {
	address := mxHost + ":" + port

	// 尝试建立连接
	conn, err := net.DialTimeout("tcp", address, v.timeout)
	if err != nil {
		return "", fmt.Errorf("cannot connect to %s: %v", address, err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return "", fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// 如果支持 STARTTLS，启用它
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         mxHost,
			InsecureSkipVerify: true, // 仅用于验证，不发送实际邮件
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			// STARTTLS 失败，继续尝试不加密的连接
		}
	}

	// HELO 命令
	if err := client.Hello("example.com"); err != nil {
		return "", fmt.Errorf("HELO failed: %v", err)
	}

	// MAIL FROM 命令
	if err := client.Mail(v.fromEmail); err != nil {
		return "", fmt.Errorf("MAIL FROM failed: %v", err)
	}

	// RCPT TO 命令 - 这是关键步骤
	if err := client.Rcpt(email); err != nil {
		// 检查错误类型
		errStr := err.Error()
		if strings.Contains(errStr, "550") ||
			strings.Contains(errStr, "551") ||
			strings.Contains(errStr, "553") ||
			strings.Contains(errStr, "User unknown") ||
			strings.Contains(errStr, "does not exist") {
			// 邮箱不存在
			client.Quit()
			return "dead", nil
		}
		// 其他错误（可能是临时错误或服务器配置问题）
		return "", fmt.Errorf("RCPT TO failed: %v", err)
	}

	// QUIT 命令
	client.Quit()

	// 邮箱存在且可用
	return "live", nil
}

// VerifyEmailBatch 批量验证邮箱
func (v *SMTPVerifier) VerifyEmailBatch(emails []string) map[string]string {
	results := make(map[string]string)

	for _, email := range emails {
		status, err := v.VerifyEmail(email)
		if err != nil {
			// 记录错误但继续处理
			fmt.Printf("Error verifying %s: %v\n", email, err)
		}
		results[email] = status

		// 添加延迟避免被限流
		time.Sleep(500 * time.Millisecond)
	}

	return results
}

// VerifyEmailQuick 快速验证（仅检查 MX 记录）
func (v *SMTPVerifier) VerifyEmailQuick(email string) (string, error) {
	// 1. 验证邮箱格式
	if !strings.Contains(email, "@") {
		return "dead", fmt.Errorf("invalid email format")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "dead", fmt.Errorf("invalid email format")
	}
	domain := parts[1]

	// 2. 查询 MX 记录
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return "dead", fmt.Errorf("no MX records found for domain: %s", domain)
	}

	// 有 MX 记录，但不确定邮箱是否真实存在
	return "verify", nil
}
