#!/bin/bash

while true; do
  _issue_count=$(bd count --status open)
  echo "Issues in progress: ${_issue_count}"
  if [ "$_issue_count" -eq 0 ]; then
    echo "No issues in progress. Exiting..."
    break
  fi
  echo ""
  echo "Starting AI-assisted development cycle"
  echo ""
  cat PROMPT.txt | claude --dangerously-skip-permissions -p --model haiku
  echo ""
done