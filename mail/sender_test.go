package mail

import (
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from SimpleBank App</p>
`
	to := []string{"20200255@stud.fci-cu.edu.eg"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
