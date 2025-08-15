#!/bin/bash

echo "🧪 Testing FinHub Picklist System"
echo "=================================="

# Check if backend is running
echo "📡 Checking backend status..."
if curl -s http://localhost:8080/api/auth/login > /dev/null 2>&1; then
    echo "✅ Backend is running"
else
    echo "❌ Backend is not running. Please start it first with:"
    echo "   cd backend && make run"
    exit 1
fi

# Check if we have a test user
echo ""
echo "🔐 Setting up authentication..."
echo "   Note: You need to create a test user first. Run:"
echo "   cd backend && ./mkuser.sh"
echo ""
echo "   Or manually create a user in the database."
echo ""
echo "   For now, we'll test the endpoints that don't require auth:"
echo "   - POST /api/auth/login"
echo "   - POST /api/auth/register"

# Test authentication endpoints
echo ""
echo "🔍 Testing authentication endpoints..."

echo "📝 Testing registration endpoint..."
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123","firstName":"Test","lastName":"User"}')
if echo "$REGISTER_RESPONSE" | grep -q "token"; then
    echo "✅ Registration endpoint working"
    echo "   User created successfully"
elif echo "$REGISTER_RESPONSE" | grep -q "already exists"; then
    echo "✅ Registration endpoint working"
    echo "   User already exists (this is fine)"
else
    echo "❌ Registration endpoint failed"
    echo "   Response: $REGISTER_RESPONSE"
fi

echo "🔑 Testing login endpoint..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}')
if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo "✅ Login endpoint working"
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null || echo "")
    if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
        echo "   Token received successfully"
        echo ""
        echo "🔍 Now testing picklist endpoints with authentication..."
        
        echo "📋 Testing industries endpoint..."
        INDUSTRIES_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/picklists/industries)
        if echo "$INDUSTRIES_RESPONSE" | grep -q "items"; then
            echo "✅ Industries endpoint working"
            echo "   Found $(echo "$INDUSTRIES_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown") industries"
        else
            echo "❌ Industries endpoint failed"
            echo "   Response: $INDUSTRIES_RESPONSE"
        fi

        echo "📊 Testing company sizes endpoint..."
        SIZES_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/picklists/companysizes)
        if echo "$SIZES_RESPONSE" | grep -q "items"; then
            echo "✅ Company sizes endpoint working"
            echo "   Found $(echo "$SIZES_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown") company sizes"
        else
            echo "❌ Company sizes endpoint failed"
            echo "   Response: $SIZES_RESPONSE"
        fi

        echo "🎯 Testing search endpoint..."
        SEARCH_RESPONSE=$(curl -s -X POST http://localhost:8080/api/picklists/search \
          -H "Authorization: Bearer $TOKEN" \
          -H "Content-Type: application/json" \
          -d '{"query":"tech","limit":10,"offset":0,"entityType":"industry"}')
        if echo "$SEARCH_RESPONSE" | grep -q "items"; then
            echo "✅ Search endpoint working"
            echo "   Search results: $(echo "$SEARCH_RESPONSE" | jq '.items | length' 2>/dev/null || echo "unknown") items"
        else
            echo "❌ Search endpoint failed"
            echo "   Response: $SEARCH_RESPONSE"
        fi
    else
        echo "❌ Failed to extract token from response"
    fi
else
    echo "❌ Login endpoint failed"
    echo "   Response: $LOGIN_RESPONSE"
fi

echo ""
echo "🎉 Picklist system test completed!"
echo ""
echo "💡 To test the frontend:"
echo "   1. Start the frontend: cd frontend && npm run dev"
echo "   2. Open http://localhost:3000"
echo "   3. Try creating a company with the new picklist dropdowns"
echo ""
echo "📚 For more information, see PICKLIST_README.md" 