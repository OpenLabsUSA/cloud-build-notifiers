package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"
	cbpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func TestWriteMessage(t *testing.T) {
	n := new(slackNotifier)
	b := &cbpb.Build{
		ProjectId: "my-project-id",
		Id:        "some-build-id",
		Status:    cbpb.Build_SUCCESS,
		LogUrl:    "https://some.example.com/log/url?foo=bar",
		Substitutions: map[string]string{
			"REPO_NAME": "repo-name",
			"SHORT_SHA": "abc123",
		},
		Steps: []*cbpb.BuildStep{
			{Id: "step 1", Status: cbpb.Build_SUCCESS},
			{Id: "step 2", Status: cbpb.Build_FAILURE},
		},
	}

	got, err := n.writeMessage(b)
	if err != nil {
		t.Fatalf("writeMessage failed: %v", err)
	}

	want := &slack.WebhookMessage{
		Attachments: []slack.Attachment{{
			Text:  "status: **SUCCESS**\nrepo: `repo-name`\ncommit: `abc123`\nstep: **step 2**",
			Color: "good",
			Actions: []slack.AttachmentAction{{
				Text: "View Logs",
				Type: "button",
				URL:  "https://some.example.com/log/url?foo=bar&utm_campaign=google-cloud-build-notifiers&utm_medium=chat&utm_source=google-cloud-build",
			}},
		}},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("writeMessage got unexpected diff: %s", diff)
	}
}
