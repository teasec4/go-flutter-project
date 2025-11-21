#!/bin/bash

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Bank API Server Test Script ===${NC}\n"

# Build the server
echo -e "${YELLOW}Building server...${NC}"
cd cmd/api
if ! go build -o api .; then
    echo -e "${RED}Build failed${NC}"
    exit 1
fi
cd ../..
echo -e "${GREEN}✓ Build successful${NC}\n"

# Start the server
echo -e "${YELLOW}Starting server...${NC}"
./cmd/api/api > /tmp/server.log 2>&1 &
SERVER_PID=$!
sleep 2
echo -e "${GREEN}✓ Server started (PID: $SERVER_PID)${NC}\n"

# Function to make requests
test_request() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_code=$5
    
    echo -e "${BLUE}→ Testing: $name${NC}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "http://localhost:8080$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "http://localhost:8080$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" == "$expected_code" ]; then
        echo -e "${GREEN}✓ HTTP $http_code (expected $expected_code)${NC}"
        echo "  Response: $body" | jq . 2>/dev/null || echo "  Response: $body"
    else
        echo -e "${RED}✗ HTTP $http_code (expected $expected_code)${NC}"
        echo "  Response: $body"
    fi
    echo ""
}

# Run tests
echo -e "${BLUE}=== API Tests ===${NC}\n"

test_request "Get Account 1 Balance" "GET" "/account?id=1" "" "200"
test_request "Deposit 500 to Account 1" "POST" "/account/deposit" '{"accountId":"1","amount":500}' "200"
test_request "Get Updated Balance" "GET" "/account?id=1" "" "200"
test_request "Withdraw 200 from Account 1" "POST" "/account/withdraw" '{"accountId":"1","amount":200}' "200"
test_request "Get Final Balance" "GET" "/account?id=1" "" "200"
test_request "Get Account 2 Balance" "GET" "/account?id=2" "" "200"

echo -e "${BLUE}=== Error Tests ===${NC}\n"

test_request "Missing ID Parameter" "GET" "/account" "" "400"
test_request "Non-existent Account" "GET" "/account?id=999" "" "404"
test_request "Withdraw More Than Balance" "POST" "/account/withdraw" '{"accountId":"1","amount":50000}' "400"
test_request "Negative Deposit" "POST" "/account/deposit" '{"accountId":"1","amount":-100}' "400"

# Show server logs
echo -e "${BLUE}=== Server Logs ===${NC}\n"
echo "Last 20 lines of server logs:"
tail -20 /tmp/server.log | sed 's/^/  /'

# Graceful shutdown test
echo -e "\n${YELLOW}Testing graceful shutdown...${NC}"
echo "Sending SIGINT to server..."
kill -INT $SERVER_PID

# Wait a bit for clean shutdown
sleep 2

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${GREEN}✓ Server shut down gracefully${NC}"
else
    echo -e "${YELLOW}⚠ Server still running, killing forcefully...${NC}"
    kill -9 $SERVER_PID 2>/dev/null
fi

echo -e "\n${GREEN}=== All Tests Completed ===${NC}"
