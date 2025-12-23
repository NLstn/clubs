package notifications

import (
	"fmt"
	"os"

	"github.com/NLstn/civo/azure/acs"
	frontend "github.com/NLstn/civo/tools"
)

// SendMemberAddedNotification sends both email and in-app notifications for member addition
// This function should be called from the models package after creating in-app notification
func SendMemberAddedNotification(userEmail, clubID string, clubName string) error {
	return sendMemberAddedEmail(userEmail, clubID, clubName)
}

// SendMemberAddedEmailIfEnabled sends email notification if user preferences allow it
// This is a separate function to avoid circular imports
func SendMemberAddedEmailIfEnabled(userEmail, clubID string, clubName string, emailEnabled bool) error {
	if emailEnabled {
		return sendMemberAddedEmail(userEmail, clubID, clubName)
	}
	return nil
}

// sendMemberAddedEmail sends the email notification for member addition
func sendMemberAddedEmail(userMail string, clubID string, clubName string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	subject := "You have been added to a club"
	clubLink := frontend.MakeClubLink(clubID)

	plainText := fmt.Sprintf("Hello,\n\nYou have been added to the club %s as a member.\n\nVisit the club page at: %s\n\nBest regards,\nThe Clubs Team", clubName, clubLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>You have been added to the club <strong>%s</strong> as a member.</p><p>Visit the club page <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", clubName, clubLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}

// sendEventCreatedEmail sends the email notification for event creation
func sendEventCreatedEmail(userMail, clubID, eventID, eventTitle string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	subject := "New event: " + eventTitle
	eventLink := frontend.MakeEventLink(clubID, eventID)

	plainText := fmt.Sprintf("Hello,\n\nA new event '%s' has been created.\n\nView the event at: %s\n\nBest regards,\nThe Clubs Team", eventTitle, eventLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>A new event <strong>%s</strong> has been created.</p><p>View the event <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", eventTitle, eventLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}

// sendFineAssignedEmail sends the email notification for fine assignment
func sendFineAssignedEmail(userMail, clubID, fineID string, fineAmount float64, reason string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	subject := "Fine assigned"
	fineLink := frontend.MakeFineLink(clubID, fineID)

	plainText := fmt.Sprintf("Hello,\n\nYou have been assigned a fine of €%.2f for: %s\n\nView your fine at: %s\n\nBest regards,\nThe Clubs Team", fineAmount, reason, fineLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>You have been assigned a fine of <strong>€%.2f</strong> for: %s</p><p>View your fine <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", fineAmount, reason, fineLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}

// sendNewsCreatedEmail sends the email notification for news creation
func sendNewsCreatedEmail(userMail, clubID, newsTitle string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	subject := "New news: " + newsTitle
	clubLink := frontend.MakeClubLink(clubID)

	plainText := fmt.Sprintf("Hello,\n\nA new news post '%s' has been published.\n\nView the news at: %s\n\nBest regards,\nThe Clubs Team", newsTitle, clubLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>A new news post <strong>%s</strong> has been published.</p><p>View the news <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", newsTitle, clubLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}

// SendEventCreatedEmailIfEnabled sends email notification for new events if enabled
func SendEventCreatedEmailIfEnabled(userEmail, clubID, eventID string, eventTitle string, emailEnabled bool) error {
	if emailEnabled {
		return sendEventCreatedEmail(userEmail, clubID, eventID, eventTitle)
	}
	return nil
}

// SendFineAssignedEmailIfEnabled sends email notification for fine assignments if enabled
func SendFineAssignedEmailIfEnabled(userEmail, clubID, fineID string, fineAmount float64, reason string, emailEnabled bool) error {
	if emailEnabled {
		return sendFineAssignedEmail(userEmail, clubID, fineID, fineAmount, reason)
	}
	return nil
}

// SendNewsCreatedEmailIfEnabled sends email notification for new news posts if enabled
func SendNewsCreatedEmailIfEnabled(userEmail, clubID string, newsTitle string, emailEnabled bool) error {
	if emailEnabled {
		return sendNewsCreatedEmail(userEmail, clubID, newsTitle)
	}
	return nil
}

// SendRoleChangedNotification sends email notification for role changes
func SendRoleChangedNotification(userEmail, clubID, clubName, oldRole, newRole string) error {
	return sendRoleChangedEmail(userEmail, clubID, clubName, oldRole, newRole)
}

// SendRoleChangedEmailIfEnabled sends email notification for role changes if enabled
func SendRoleChangedEmailIfEnabled(userEmail, clubID, clubName, oldRole, newRole string, emailEnabled bool) error {
	if emailEnabled {
		return sendRoleChangedEmail(userEmail, clubID, clubName, oldRole, newRole)
	}
	return nil
}

// sendRoleChangedEmail sends the email notification for role changes
func sendRoleChangedEmail(userMail, clubID, clubName, oldRole, newRole string) error {
	// Skip Azure Communication Services email calls in test environment
	if os.Getenv("GO_ENV") == "test" {
		return nil
	}

	subject := "Role updated in " + clubName
	clubLink := frontend.MakeClubLink(clubID)

	plainText := fmt.Sprintf("Hello,\n\nYour role in the club %s has been changed from %s to %s.\n\nVisit the club page at: %s\n\nBest regards,\nThe Clubs Team", clubName, oldRole, newRole, clubLink)
	htmlContent := fmt.Sprintf("<p>Hello,</p><p>Your role in the club <strong>%s</strong> has been changed from <strong>%s</strong> to <strong>%s</strong>.</p><p>Visit the club page <a href=\"%s\">here</a>.</p><p>Best regards,<br>The Clubs Team</p>", clubName, oldRole, newRole, clubLink)

	recipients := []acs.Recipient{
		{Address: userMail},
	}

	return acs.SendMail(recipients, subject, plainText, htmlContent)
}
