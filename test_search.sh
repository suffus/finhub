#!/bin/bash

echo "🔍 Testing Picklist Search Functionality"
echo "========================================"

# Check if backend is running
echo "📡 Checking backend status..."
if curl -s http://localhost:8080/api/auth/login > /dev/null 2>&1; then
    echo "✅ Backend is running"
else
    echo "❌ Backend is not running. Please start it first with:"
    echo "   cd backend && make run"
    exit 1
fi

echo ""
echo "🔐 Testing authentication and search..."

# Register a test user (if not exists)
echo "📝 Creating test user..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"searchtest@example.com","password":"testpass123","firstName":"Search","lastName":"Test"}')

if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    echo "✅ User created successfully"
elif echo "$REGISTER_RESPONSE" | grep -q "already exists"; then
    echo "✅ User already exists (this is fine)"
else
    echo "❌ Failed to create user"
    echo "   Response: $REGISTER_RESPONSE"
    exit 1
fi

# Login to get token
echo "🔑 Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"searchtest@example.com","password":"testpass123"}')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo "✅ Login successful"
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null || echo "")
    
    if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
        echo "✅ Token received"
        echo ""
        echo "🔍 Testing search functionality..."
        
        # Test 1: Search for industries with "tech"
        echo "📋 Testing industry search for 'tech'..."
        SEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/picklists/search \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{"query":"tech","limit":10,"offset":0,"entityType":"industry"}')
        
        if echo "$SEARCH_RESPONSE" | grep -q "items"; then
            echo "✅ Industry search working"
            ITEM_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown")
            echo "   Found $ITEM_COUNT items"
        else
            echo "❌ Industry search failed"
            echo "   Response: $SEARCH_RESPONSE"
        fi
        
        # Test 2: Search for company sizes with "startup"
        echo "📊 Testing company size search for 'startup'..."
        SEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/picklists/search \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{"query":"startup","limit":10,"offset":0,"entityType":"companysize"}')
        
        if echo "$SEARCH_RESPONSE" | grep -q "items"; then
            echo "✅ Company size search working"
            ITEM_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown")
            echo "   Found $ITEM_COUNT items"
        else
            echo "❌ Company size search failed"
            echo "   Response: $SEARCH_RESPONSE"
        fi
        
        # Test 3: Empty search query
        echo "🔍 Testing empty search query..."
        SEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/picklists/search \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{"query":"","limit":10,"offset":0,"entityType":"industry"}')
        
        if echo "$SEARCH_RESPONSE" | grep -q "items"; then
            echo "✅ Empty search working"
            ITEM_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown")
            echo "   Found $ITEM_COUNT items"
        else
            echo "❌ Empty search failed"
            echo "   Response: $SEARCH_RESPONSE"
        fi
        
        # Test 4: Invalid entity type
        echo "🚫 Testing invalid entity type..."
        SEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/picklists/search \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{"query":"test","limit":10,"offset":0,"entityType":"invalid_type"}')
        
        if echo "$SEARCH_RESPONSE" | grep -q "Invalid entity type"; then
            echo "✅ Invalid entity type handling working"
        else
            echo "❌ Invalid entity type handling failed"
            echo "   Response: $SEARCH_RESPONSE"
        fi
        
    else
        echo "❌ Failed to extract token"
        exit 1
    fi
else
    echo "❌ Login failed"
    echo "   Response: $LOGIN_RESPONSE"
    exit 1
fi

echo ""
echo "🎉 Search functionality test completed!"
echo ""
echo "💡 If all tests passed, the search functionality is working correctly."
echo "   If any failed, check the backend logs for detailed error messages." 