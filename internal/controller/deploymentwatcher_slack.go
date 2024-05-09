package controller

import (
	"github.com/slack-go/slack"
	"strings"
)

func getTag(version string) string {
	splitVersion := strings.Split(version, ":")
	return splitVersion[len(splitVersion)-1]
}

func notify(namespace, name, version, slackToken, slackChannel string) error {

	tag := getTag(version)

	headerSection := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", ":rocket: New deployment is in progress!", false, false),
		nil,
		nil,
	)

	fields := make([]*slack.TextBlockObject, 0)
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", "*Namespace/Deployment*", false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", "*Version*", false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", namespace+"/"+name, false, false))
	fields = append(fields, slack.NewTextBlockObject("mrkdwn", tag, false, false))

	deploymentSection := slack.NewSectionBlock(
		nil,
		fields,
		nil,
	)

	attachment := slack.Attachment{
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{headerSection, deploymentSection},
		},
		Color: "#36a64f",
	}

	message := slack.MsgOptionAttachments(attachment)
	api := slack.New(slackToken)
	_, _, err := api.PostMessage(slackChannel, message)
	if err != nil {
		return err
	}
	return nil
}
