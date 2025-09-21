#!/bin/bash

# Make sure to run the database seeder first to ensure the test user exists.
# You will need curl and jq installed to run this script.

# --- Configuration ---
API_URL="http://localhost:8080/query"
EMAIL="test@example.com"
PASSWORD="password123"

# --- 1. Log in and get JWT token ---
echo "Logging in as $EMAIL..."
LOGIN_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" \
  -d "{\"query\":\"mutation { login(email: \"$EMAIL\", password: \"$PASSWORD\") { token } }\"}" \
  $API_URL)

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.login.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "Failed to get JWT token. Response was:"
  echo $LOGIN_RESPONSE
  exit 1
fi

echo "Successfully logged in. Token received."
echo "----------------------------------------"

# --- 2. Send a burst of requests ---
echo "Sending 5 requests in rapid succession..."

for i in {1..5}
do
  echo -n "Request $i: "
  curl -s -o /dev/null -w "%{http_code}\n" -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"query":"query { me { id email } }"}' \
    $API_URL
done

echo "----------------------------------------"
echo "Test complete."
