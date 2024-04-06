#!/bin/bash

SPIN_VARIABLE_ALLOWED_CHANNEL=slack-channel-id \
SPIN_VARIABLE_TRIGGER_ON_EMOJI_CODE=shipit \
SPIN_VARIABLE_SLACK_TOKEN=xoxb-slack-token \
SPIN_VARIABLE_SLACK_SIGNING_SECRET=slack-signing-token \
spin up \
	-f ghcr.io/rajatjindal/behind-the-scenes:v0.1.0 \
	--disable-pooling \
	--state-dir ./.spin \
	--listen 0.0.0.0:3000 \
	--key-value kv-credentials=kvexplorer-username:kvexplorer-password