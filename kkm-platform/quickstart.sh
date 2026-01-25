#!/bin/bash

# KKM Platform Quick Start Script
# This script helps you get started with the KKM Platform

set -e

echo "🚀 KKM Platform Quick Start"
echo "=============================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running from correct directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: package.json not found"
    echo "Please run this script from the project root directory"
    exit 1
fi

# 1. Check prerequisites
echo -e "${BLUE}1. Checking prerequisites...${NC}"

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed"
    exit 1
fi

if ! command -v pnpm &> /dev/null; then
    echo "⚠️  pnpm is not installed, trying npm..."
    PACKAGE_MANAGER="npm"
else
    PACKAGE_MANAGER="pnpm"
fi

NODE_VERSION=$(node --version)
echo -e "${GREEN}✓ Node.js ${NODE_VERSION} found${NC}"
echo -e "${GREEN}✓ Using package manager: ${PACKAGE_MANAGER}${NC}"
echo ""

# 2. Install dependencies
echo -e "${BLUE}2. Installing dependencies...${NC}"
$PACKAGE_MANAGER install
echo -e "${GREEN}✓ Dependencies installed${NC}"
echo ""

# 3. Setup environment
echo -e "${BLUE}3. Setting up environment...${NC}"
if [ ! -f ".env.local" ]; then
    cat > .env.local << EOF
# Development environment variables
NEXT_PUBLIC_API_URL=http://localhost:3000
NEXT_PUBLIC_WS_URL=ws://localhost:3000
EOF
    echo -e "${GREEN}✓ .env.local created${NC}"
else
    echo -e "${YELLOW}⚠️  .env.local already exists${NC}"
fi
echo ""

# 4. Run build check
echo -e "${BLUE}4. Running production build check...${NC}"
npm run build > /dev/null 2>&1
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# 5. Display next steps
echo -e "${BLUE}5. Next steps:${NC}"
echo ""
echo -e "${GREEN}Start development server:${NC}"
echo "  npm run dev"
echo ""
echo -e "${GREEN}Run tests:${NC}"
echo "  npm run test"
echo ""
echo -e "${GREEN}View available scripts:${NC}"
echo "  npm run"
echo ""

# 6. Demo pages info
echo -e "${BLUE}Demo Pages:${NC}"
echo "  📡 WebSocket Demo:         http://localhost:3000/realtime-demo"
echo "  📊 Audit Demo:             http://localhost:3000/audit-demo"
echo "  🌍 i18n Demo:              http://localhost:3000/ru/i18n-demo"
echo "  🌙 Dark Mode Demo:         http://localhost:3000/ru/theme-demo"
echo ""

# 7. Documentation links
echo -e "${BLUE}Documentation:${NC}"
echo "  📚 Full Index:             DOCUMENTATION_FULL_INDEX.md"
echo "  🔄 WebSocket:              WEBSOCKET_REALTIME.md"
echo "  📝 Audit Logging:          AUDIT_LOGGING.md"
echo "  🌍 Internationalization:   INTERNATIONALIZATION.md"
echo "  🌙 Dark Mode:              DARK_MODE.md"
echo ""

echo -e "${GREEN}✅ Setup complete!${NC}"
echo ""
echo "Run 'npm run dev' to start development server"
