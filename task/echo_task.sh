#!/bin/bash
MESSAGE=${1:-"No message provided"}
echo "Echo Task Result:"
echo "================"
echo "Message: $MESSAGE"
echo "Received at: $(date)"
echo "From: $(hostname)"
