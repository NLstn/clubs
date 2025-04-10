package notifications

import (
	"fmt"

	"github.com/NLstn/clubs/azure/acs"
	frontend "github.com/NLstn/clubs/tools"
)

func SendMemberAddedNotification(userMail string, clubName string, clubID string) error {
	subject := "You have been added to a club"
	clubLink := frontend.MakeClubLink(clubID)

	plainText := fmt.Sprintf("Hello,\n\nYou have been added to the club %s as a member.\n\nVisit the club page at: %s\n\nBest regards,\nThe Clubs Team", clubName, clubLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>You have been added to the club <strong>%s</strong> as a member.</p><p>Visit the club page <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", clubName, clubLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}
