#!/bin/bash

PROMPT="
read the @docs/specs/SPEC-tags.md
get the current available work 'what should I work on'
Pick item 1
work on the task until complete
create issues for any issues that you find
run 'vibeguard check' and fix any failures raised
log your findings
commit your changes
talk about what happened
"
ITERATION=1
while true; do
  _issue_open_count=$(bd count --status open)
  _issue_in_progress_count=$(bd count --status "in_progress")
  echo "Issues open: ${_issue_open_count}"
  echo "Issues in progress: ${_issue_in_progress_count}"
  if [ "$_issue_open_count" -eq 0 ] && [ "$_issue_in_progress_count" -eq 0 ]; then
    echo "No issues open or in progress. Exiting..."
    break
  fi
  echo ""
  echo "[Iteration ${ITERATION}] Starting AI-assisted development cycle"
  echo ""
  echo "${PROMPT}" | claude --dangerously-skip-permissions -p --model haiku
  echo ""
  ITERATION=$((ITERATION + 1))
done