package notifications

import (
	"fmt"

	"github.com/NLstn/clubs/azure/acs"
)

func SendMemberAddedNotification(userMail string, clubName string) error {
	subject := "You have been added to a club"
	plainText := fmt.Sprintf("Hello,\n\nYou have been added to the club %s as a member.\n\nBest regards,\nThe Clubs Team", clubName)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>You have been added to the club <strong>%s</strong> as a member.</p><p>Best regards,<br>The Clubs Team</p>", clubName)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}
