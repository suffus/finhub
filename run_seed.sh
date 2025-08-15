#!/bin/bash

echo "🌱 Building and Running Enhanced Database Seeder"
echo "================================================"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the backend directory"
    exit 1
fi

echo "🔨 Building the seeder..."
if go build -o seed cmd/seed/main.go; then
    echo "✅ Seeder built successfully"
else
    echo "❌ Failed to build seeder"
    exit 1
fi

echo ""
echo "🌱 Running the seeder..."
echo "   This will populate the database with comprehensive sample data including:"
echo "   - Industries and Company Sizes"
echo "   - Lead Statuses and Temperatures"
echo "   - Marketing Source Types and Sources"
echo "   - Marketing Asset Types"
echo "   - Communication Types"
echo "   - Task Types"
echo "   - Address, Email, and Phone Types"
echo "   - Territory Types and Sample Territories"
echo "   - Social Media Types"
echo "   - Sales Pipeline and Stages"
echo ""

./seed

echo ""
echo "🎉 Seeding completed!"
echo ""
echo "💡 You can now:"
echo "   1. Test the picklists in the frontend"
echo "   2. Create companies with industry and size data"
echo "   3. Use the comprehensive CRM system"
echo ""
echo "🔍 To verify the data, check the database or test the API endpoints" 