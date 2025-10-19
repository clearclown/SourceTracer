#!/bin/bash
# E2E Test Script for SourceTracer API

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="${API_URL:-http://localhost:8080/api/v1}"
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test helper functions
test_start() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "${YELLOW}[TEST $TOTAL_TESTS]${NC} $1"
}

test_pass() {
    PASSED_TESTS=$((PASSED_TESTS + 1))
    echo -e "${GREEN}✓ PASSED${NC}"
    echo ""
}

test_fail() {
    FAILED_TESTS=$((FAILED_TESTS + 1))
    echo -e "${RED}✗ FAILED${NC}: $1"
    echo ""
}

# Test 1: Health Check
test_start "Health check endpoint"
RESPONSE=$(curl -s "$API_URL/health")
if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    test_pass
else
    test_fail "Health check returned unexpected response"
fi

# Test 2: Analyze simple text
test_start "Analyze simple opinion text"
RESPONSE=$(curl -s -X POST "$API_URL/analyze" \
    -H "Content-Type: application/json" \
    -d '{"text":"Python is the best programming language.","options":{"include_evidences":false}}')

if echo "$RESPONSE" | grep -q '"success":true'; then
    if echo "$RESPONSE" | grep -q '"claims"'; then
        test_pass
    else
        test_fail "Response missing claims field"
    fi
else
    test_fail "Analyze request failed"
fi

# Test 3: Analyze multi-claim text
test_start "Analyze multi-claim text"
RESPONSE=$(curl -s -X POST "$API_URL/analyze" \
    -H "Content-Type: application/json" \
    -d '{"text":"Climate change is accelerating. Python is great. The Earth is round.","options":{"include_evidences":false,"max_claims":10}}')

CLAIM_COUNT=$(echo "$RESPONSE" | jq '.data.claims | length' 2>/dev/null || echo "0")
if [ "$CLAIM_COUNT" -ge 2 ]; then
    test_pass
else
    test_fail "Expected at least 2 claims, got $CLAIM_COUNT"
fi

# Test 4: Empty text handling
test_start "Handle empty text gracefully"
RESPONSE=$(curl -s -X POST "$API_URL/analyze" \
    -H "Content-Type: application/json" \
    -d '{"text":"","options":{}}')

if echo "$RESPONSE" | grep -q '"success":false'; then
    test_pass
else
    test_fail "Should reject empty text"
fi

# Test 5: History endpoint
test_start "History endpoint returns list"
RESPONSE=$(curl -s "$API_URL/history")

if echo "$RESPONSE" | grep -q '"success":true'; then
    if echo "$RESPONSE" | grep -q '"analyses"'; then
        test_pass
    else
        test_fail "History response missing analyses field"
    fi
else
    test_fail "History request failed"
fi

# Test 6: History with pagination
test_start "History endpoint with pagination"
RESPONSE=$(curl -s "$API_URL/history?limit=5&offset=0")

LIMIT=$(echo "$RESPONSE" | jq '.data.limit' 2>/dev/null || echo "0")
if [ "$LIMIT" -eq 5 ]; then
    test_pass
else
    test_fail "Pagination not working correctly"
fi

# Test 7: Get specific analysis (should fail with mock DB)
test_start "Get non-existent analysis returns 404"
RESPONSE=$(curl -s "$API_URL/history/non-existent-id")

if echo "$RESPONSE" | grep -q '"success":false'; then
    test_pass
else
    test_fail "Should return error for non-existent ID"
fi

# Test 8: Invalid JSON handling
test_start "Handle invalid JSON gracefully"
RESPONSE=$(curl -s -X POST "$API_URL/analyze" \
    -H "Content-Type: application/json" \
    -d 'invalid json')

if echo "$RESPONSE" | grep -q '"success":false'; then
    test_pass
else
    test_fail "Should reject invalid JSON"
fi

# Summary
echo "========================================"
echo -e "${GREEN}PASSED: $PASSED_TESTS${NC}"
echo -e "${RED}FAILED: $FAILED_TESTS${NC}"
echo "TOTAL:  $TOTAL_TESTS"
echo "========================================"

if [ "$FAILED_TESTS" -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
