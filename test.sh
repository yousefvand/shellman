#!/usr/bin/env bash

# Usage: progressBar "message" currentStep totalSteps
function progressBar() {
  local   bar='████████████████████'
  local space='                    '
  local wheel=('\' '|' '/' '-')

  local msg="${1}"
  local current=${2}
  local total=${3}
  local wheelIndex=$((current % 4))
  local position=$((20 * current / total))

  echo -ne "\r|${bar:0:$position}${space:$position:20}| ${wheel[wheelIndex]} $(($position*5))% [ ${msg} ] "
}

# Usage: progressBar "message" currentStep totalSteps
totalSteps=100

for ((currentStep=1; currentStep <= totalSteps; currentStep++)); do
  sleep 0.1 # simulating one step of job
  progressBar "Installing foo..." ${currentStep} ${totalSteps}
done

echo
