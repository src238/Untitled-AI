# FinAI Work Log
**Date**: January 31, 2026
**Session Start**: 20:32 UTC

## Summary
Successfully started both backend and frontend servers for the FinAI application. Fixed notice board to display newest alerts at the top. **Completely overhauled the AI system** to automatically check all purchases from the past week for cheaper alternatives while preserving quality, posting recommendations to the notice board every 30 seconds. Made transaction history scrollable and **loaded all 60 mock transactions** from the data file into the UI via new API endpoint. **Added clickable purchase links** to every alternative recommendation so users can immediately buy the cheaper products. **Implemented $5 minimum savings filter** to only show meaningful alternatives. **Enhanced format** to clearly show original purchase, alternative product, exact savings, and total prices. **Fixed transaction arrows** to correctly indicate money in (↑) vs money out (↓).

---

## Tasks Completed

### 1. Backend Server Startup ✅ COMPLETED
**Location**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI`
**Command**: `go run main.go`
**Status**: Running in background (Process ID: b2597b4)
**Started**: 20:32:55

#### Backend Features Initialized:
- ✅ Liminal API configured
- ✅ Added 8 Liminal banking tools (get_transactions disabled - using mock data)
- ✅ Added custom spending analyzer tool
- ✅ Added AI-powered product analyzer tool
- ✅ Added mock transaction reader tool
- ✅ Added product search tool for finding alternatives
- ✅ Added alert notification tools (post & read) for AI insights
- ✅ Started AI background analysis loop (runs every 30 seconds)

#### Backend Endpoints:
- 📡 WebSocket: `ws://localhost:8080/ws`
- 🔔 Alerts API: `http://localhost:8080/api/alerts`
- 💚 Health check: `http://localhost:8080/health`

**Output Log**: `/tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b2597b4.output`

---

### 2. Frontend Server Startup ✅ COMPLETED
**Location**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend`
**Command**: `npm run dev`
**Status**: Running in background (Process ID: b5be003)
**Started**: ~20:33:00

#### Frontend Configuration:
- 🌐 Local URL: `http://localhost:5173/`
- ⚡ Built with: Vite v5.4.21
- 📦 Dependencies: 115 node_modules installed
- ⏱️ Ready in: 166ms

**Output Log**: `/tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b5be003.output`

---

### 3. Notice Board Order Fix ✅ COMPLETED
**File Modified**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/main.tsx`
**Time**: 20:45:00
**Status**: Applied and hot-reloaded

#### Changes Made:
- Modified alert display logic to show newest notices at the top
- Changed line 78: Added `.reverse()` to `mappedAlerts` array before setting state
- Ensures consistency: both API-fetched alerts and WebSocket real-time alerts now appear with newest first
- Frontend automatically reloaded via Vite HMR (Hot Module Replacement)

#### Technical Details:
- **File**: `frontend/main.tsx:78`
- **Change**: `setAlerts(mappedAlerts)` → `setAlerts(mappedAlerts.reverse())`
- **Result**: AI Insights sidebar now displays newest alerts at top, pushing older ones down
- **Behavior**: New alerts arrive at top every ~30 seconds from background AI analysis

---

## Project Structure
```
finAI/
├── main.go                      # Go backend server with AI agent
├── frontend/                    # React chat interface
│   ├── main.tsx                 # App entry point (9,058 bytes)
│   ├── index.html               # Landing page
│   ├── styles.css               # Styling (10,213 bytes)
│   ├── vite.config.ts           # Vite configuration
│   └── node_modules/            # 115 packages installed
├── go.mod                       # Go dependencies
├── go.sum                       # Go dependency checksums
├── .env                         # Environment configuration (1,709 bytes)
├── mock_transactions.txt        # Mock transaction data (6,307 bytes)
├── test_product_analyzer.go     # Test file (5,141 bytes)
├── TEST_RESULTS.md              # Test documentation
├── USAGE.md                     # Usage guide
└── README.md                    # Project documentation
```

---

## How to Access the Application

1. **Frontend Interface**: Open `http://localhost:5173/` in your browser
2. **Backend API**: WebSocket connection at `ws://localhost:8080/ws`
3. **Health Check**: Visit `http://localhost:8080/health`
4. **Alerts API**: Access at `http://localhost:8080/api/alerts`

---

## Next Steps for Users

1. Click the chat bubble in the frontend
2. Login with your email (you'll receive an OTP code)
3. Enter the code to authenticate
4. Start chatting! Try:
   - "What's my balance?"
   - "Show me my recent transactions"
   - "Analyze my spending over the last 30 days"

---

## Server Management

### Check Server Status
```bash
# Check backend logs
tail -f /tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b2597b4.output

# Check frontend logs
tail -f /tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b5be003.output
```

### Stop Servers
```bash
# Find and kill processes
lsof -ti:8080 | xargs kill -9  # Backend
lsof -ti:5173 | xargs kill -9  # Frontend
```

---

## Notes
- Backend is using mock transaction data (get_transactions from Liminal API disabled)
- AI background analysis runs every 30 seconds
- Frontend npm packages have 2 moderate severity vulnerabilities (optional to fix with `npm audit fix`)
- All environment variables configured in `.env` file

---

---

### 4. AI System Overhaul - Cheaper Alternatives Checker ✅ COMPLETED
**Files Modified**:
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/main.go`
**Time**: 21:26:15
**Status**: Deployed and actively running

#### Major Changes:
**Removed:**
- All previous generic spending analysis prompts
- Random financial insight generation
- Generic spending pattern alerts

**Replaced With:**
- Intelligent cheaper alternatives checker system
- Automatic purchase analysis from past week
- Quality-preserving product recommendations

#### New System Architecture:

**1. Product Tracking System:**
- Added `checkedProducts` map to track analyzed purchases
- Prevents duplicate analysis of same products
- Automatically resets after all products checked

**2. Smart Purchase Analysis:**
- Reads mock transactions from file
- Filters purchases from past 7 days
- Iterates through unchecked products every 30 seconds
- Tracks completion and resets after full cycle

**3. AI-Powered Alternative Finding:**
- Uses Claude AI to find cheaper alternatives
- Preserves or improves quality in recommendations
- Provides specific price comparisons
- Includes key benefits of alternatives
- Only posts legitimate savings opportunities

**4. Alert Format:**
Format: "[Product Name] for $[price] saves $[amount] - similar quality, [benefit]"

Examples from live system:
- "Google Nest Mini for $29 saves $20 - similar features, better Google Assistant integration"
- "Local café latte and muffin for $6.50 saves $1.95 - similar quality, fresher baked goods"

**5. Updated System Prompt:**
- Changed AI persona from "financial assistant" to "smart shopping assistant"
- Focus shifted to savings and value optimization
- Emphasizes quality preservation
- Explains automatic background analysis

#### Technical Implementation:

**Global State:**
```go
checkedProductsMutex sync.RWMutex
checkedProducts     map[string]bool
```

**Analysis Loop Logic:**
1. Wait 10 seconds on startup
2. Read and parse mock_transactions.txt
3. Filter transactions from past week
4. Find first unchecked product
5. Query Claude AI for cheaper alternative
6. Post recommendation to notice board
7. Wait 30 seconds, repeat

**AI Prompt Strategy:**
- Specific format requirements (25 words max)
- Emphasis on quality preservation
- Realistic 2026 pricing
- Concrete savings amounts
- Key benefit inclusion
- Skip if already optimal

#### Backend Logs Show Active Analysis:
```
21:26:25 🔍 Checking for cheaper alternatives: Echo Dot (5th Gen) Smart Speaker ($49.99)
21:26:27 💡 Posted alternative: Google Nest Mini for $29 saves $20...
21:26:57 🔍 Checking for cheaper alternatives: Grande Latte, Blueberry Muffin ($8.45)
21:26:59 💡 Posted alternative: Local café latte and muffin for $6.50 saves $1.95...
```

#### User Experience:
- Notice board displays newest alternatives at top (from previous fix)
- Alternatives appear automatically every 30 seconds
- Green success alerts for savings opportunities
- Blue info alerts for general recommendations
- No user action required - fully automatic

**Backend Process ID**: b166588
**Output Log**: `/tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b166588.output`

---

### 5. Transaction History Scroll ✅ COMPLETED
**File Modified**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/styles.css`
**Time**: 21:29:37
**Status**: Applied via hot reload

#### Changes:
- Added `max-height: 500px` to `.transactions-list` container
- Added `overflow-y: auto` to enable vertical scrolling
- Transaction history now scrollable when list exceeds 500px height
- Maintains clean UI with proper scroll behavior

**CSS Changes (Line 273-280):**
```css
.transactions-list {
  background: var(--white);
  border: 1px solid var(--light-border);
  border-radius: 12px;
  overflow: hidden;
  max-height: 500px;
  overflow-y: auto;
}
```

---

### 6. All Transactions Display ✅ COMPLETED
**Files Modified**:
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/main.go` (Backend API)
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/main.tsx` (Frontend fetch)
**Time**: 21:31:41
**Status**: Deployed and active

#### Backend Changes:
**Added New API Endpoint**: `/api/transactions`
- Reads all transactions from `mock_transactions.txt`
- Parses 60 transactions with full details
- Converts to frontend-compatible JSON format
- Returns complete transaction history with:
  - ID, amount, description, date, type, merchant
  - All amounts as negative for debits
  - Proper date formatting

**API Response Sample**:
```json
{
  "id": "tx-1",
  "amount": -49.99,
  "description": "Echo Dot (5th Gen) Smart Speaker",
  "date": "2026-01-30",
  "type": "debit",
  "merchant": "Amazon.com"
}
```

**Updated Server Logs** (Line 241):
```
💳 Transactions API: http://localhost:8080/api/transactions
```

#### Frontend Changes:
**Removed Hardcoded Data**:
- Cleared hardcoded 5 sample transactions
- Initialized `transactions` state as empty array

**Added Transaction Fetching**:
- New `useEffect` hook fetches transactions on component mount
- Fetches from `${apiBaseUrl}/api/transactions`
- Populates transaction list with all 60 real transactions
- Error handling for failed fetches

**Code Changes** (main.tsx:52-66):
```typescript
// Fetch transactions on mount
useEffect(() => {
  const fetchTransactions = async () => {
    try {
      const response = await fetch(`${apiBaseUrl}/api/transactions`)
      if (response.ok) {
        const data = await response.json()
        setTransactions(data)
      }
    } catch (error) {
      console.error('Failed to fetch transactions:', error)
    }
  }
  fetchTransactions()
}, [apiBaseUrl])
```

#### Result:
- **60 transactions** now displayed in scrollable transaction list
- All January 2026 purchases from mock data file
- Includes merchants, products, amounts, and dates
- Scrollable interface (from previous fix) handles all transactions
- Real transaction data instead of hardcoded samples

**Backend Process ID**: b95af31
**API Endpoint**: `http://localhost:8080/api/transactions`
**Transaction Count**: 60 items
**Output Log**: `/tmp/claude/-home-sholto-Documents-Hobbies-hackathons-nim-go-sdk/tasks/b95af31.output`

---

### 7. Product Links in Recommendations ✅ COMPLETED
**Files Modified**:
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/main.go` (AI prompt)
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/main.tsx` (Link rendering)
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/styles.css` (Link styling)
**Time**: 21:36:00
**Status**: Deployed and actively generating links

#### Backend Changes:
**Updated AI Prompt** (main.go:1210-1237):
- Changed format to require purchase links in recommendations
- New format: `"[Product Name] for $[price] saves $[amount] - [benefit]. Buy: [URL]"`
- Increased MaxTokens from 150 to 200 to accommodate URLs
- Added examples with real retailer links (Google Store, Adidas, Amazon, Trader Joe's)
- Specified rules: MUST include real, working purchase link from major retailers
- URLs should be direct product pages when possible

**Example Prompt Requirements**:
```
FORMAT YOUR RESPONSE EXACTLY LIKE THIS:
"[Product Name] for $[price] saves $[amount] - similar quality, [benefit]. Buy: [URL]"

EXAMPLES:
- "Google Nest Mini for $29 saves $20 - similar features, better voice recognition. Buy: https://store.google.com/product/google_nest_mini"
```

#### Frontend Changes:

**1. Added Link Parsing Function** (main.tsx:175-195):
```typescript
const renderAlertMessage = (message: string) => {
  const urlRegex = /(https?:\/\/[^\s]+)/g
  const parts = message.split(urlRegex)

  return parts.map((part, index) => {
    if (part.match(urlRegex)) {
      return (
        <a href={part} target="_blank" rel="noopener noreferrer"
           className="alert-link">
          {part}
        </a>
      )
    }
    return <span key={index}>{part}</span>
  })
}
```

**2. Updated Alert Rendering** (main.tsx:219):
- Changed from: `{alert.message}`
- Changed to: `{renderAlertMessage(alert.message)}`
- URLs now parsed and rendered as clickable anchor tags
- Links open in new tab with security attributes

**3. Added Link Styling** (styles.css:170-185):
```css
.alert-link {
  color: var(--orange);
  text-decoration: underline;
  font-weight: 500;
  transition: color 0.2s ease;
  word-break: break-all;
}

.alert-link:hover {
  color: var(--brown);
  text-decoration: none;
}
```

#### Live Examples Generated:
```
✓ "Google Nest Mini for $29 saves $20.99 - similar features, better voice recognition.
   Buy: https://store.google.com/product/google_nest_mini"

✓ "Dunkin' Latte and Blueberry Muffin for $6.29 saves $2.16 - similar quality, same taste.
   Buy: https://www.dunkindonuts.com"

✓ "Hulu (With Ads) for $7.99 saves $15.00 - similar content library, includes next-day TV.
   Buy: https://www.hulu.com/welcome"
```

#### User Experience:
- Every alternative recommendation includes a clickable purchase link
- Links displayed in orange color (brand color)
- Hover effect changes to brown
- Opens in new tab for safety
- Direct links to retailer websites
- Clean URL display with proper word breaking

**Backend Process ID**: bfbff14
**Status**: Links generating every 30 seconds with new recommendations

---

### 8. Enhanced Format + $5 Minimum Filter ✅ COMPLETED
**Files Modified**:
- `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/main.go` (AI prompt + filtering logic)
**Time**: 21:45:58
**Status**: Deployed and actively filtering

#### Changes Implemented:

**1. New Recommendation Format:**
Changed from simple format to comprehensive comparison:
```
OLD: "Google Nest Mini for $29 saves $20 - similar features. Buy: [URL]"

NEW: "You bought: Echo Dot (5th Gen) Smart Speaker ($49.99) | Alternative: Google Nest Mini ($29.99) | Save: $20.00 | Buy: [URL]"
```

**Format Components:**
- ✓ Original product name and price paid
- ✓ Alternative product name and price
- ✓ Exact savings amount
- ✓ Purchase link

**2. $5 Minimum Savings Filter:**
- AI instructed to only recommend if savings >= $5.00
- Backend parses "Save: $X.XX" from recommendations
- Validates savings amount >= 5.0 before posting
- Logs filtered items for transparency

**Filter Logic** (main.go:1260-1311):
```go
// Extract savings amount from recommendation
if strings.Contains(recommendation, "Save: $") {
    saveIdx := strings.Index(recommendation, "Save: $")
    saveStr := recommendation[saveIdx+7:] // Skip "Save: $"

    var savings float64
    if _, err := fmt.Sscanf(saveStr, "%f", &savings); err == nil {
        if savings >= 5.0 {
            shouldPost = true
        } else {
            log.Printf("✓ Savings too low: $%.2f (minimum $5.00)", savings)
        }
    }
}
```

#### Live Results:

**Items Posted (>= $5 savings):**
```
✓ Echo Dot → Google Nest Mini | Save: $20.00
✓ Netflix Premium → Apple TV+ | Save: $16.00
✓ Whole Foods Groceries → Costco Organic | Save: $37.84
✓ Uber to Airport → Public Transit Pass | Save: $22.50
✓ Target Household Items → Walmart Bundle | Save: $18.95
✓ Chipotle Bowl → Qdoba Bowl | Save: $5.68
✓ Best Buy USB-C Hub → Anker Hub | Save: $12.24
✓ CVS Vitamins → Nature Made | Save: $11.16
```

**Items Filtered (< $5 savings):**
```
✗ Starbucks Latte & Muffin ($8.45) - insufficient savings
✗ Shell Gas ($43.75) - insufficient savings
✗ iCloud+ Storage ($2.99) - insufficient savings
✗ Spotify Premium ($11.99) - insufficient savings
```

#### AI Prompt Updated:
```
CRITICAL RULES:
- MUST save at least $5.00 or respond: "Not enough savings (under $5)"
- Use exact format with " | " separators
- Show all prices with $ and two decimal places
```

**Backend Process ID**: b75ae45
**Filter Status**: Active - blocking recommendations under $5 savings

---

### 9. Transaction Arrow Fix ✅ COMPLETED
**File Modified**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/frontend/main.tsx`
**Time**: 21:46:00
**Status**: Applied via hot reload

#### Change:
- **Before**: Credit (↓) / Debit (↑) - backwards
- **After**: Credit (↑) / Debit (↓) - correct

**Code Change** (main.tsx:249):
```typescript
// OLD: {tx.type === 'credit' ? '↓' : '↑'}
// NEW: {tx.type === 'credit' ? '↑' : '↓'}
```

**Visual Logic:**
- ↑ Up arrow = Money coming IN (credit/deposit)
- ↓ Down arrow = Money going OUT (debit/purchase)

---

### 10. Notification Format Improvement ✅ COMPLETED
**File Modified**: `/home/sholto/Documents/Hobbies/hackathons/nim-go-sdk/finAI/main.go`
**Time**: 22:34:06
**Status**: Deployed and actively generating improved notifications

#### Changes Made:
**1. Removed "You bought:" Prefix:**
- Cleaner, more direct notification format
- Focuses on product comparison instead of purchase history

**2. Replaced Vertical Pipes with Dashes:**
- Changed: `" | "` → `" - "`
- More readable and visually appealing
- Better separation between information sections

**3. Updated Parsing Logic:**
- Modified savings extraction logic (line 1276)
- Changed from: `strings.IndexAny(saveStr, " |")`
- Changed to: `strings.IndexAny(saveStr, " -")`
- Ensures $5 minimum filter still works correctly

#### Format Changes:
**Before:**
```
"You bought: Echo Dot 5th Gen ($49.99) | Alternative: Google Nest Mini ($29.99) | Save: $20.00 | Buy: [URL]"
```

**After:**
```
"Echo Dot 5th Gen ($49.99) - Alternative: Google Nest Mini ($29.99) - Save: $20.00 - Buy: [URL]"
```

#### Code Changes:
**AI Prompt Format (lines 1222, 1225-1227, 1232):**
- Updated format template and all examples
- Added note: "no vertical pipes" in critical rules
- Maintains all information: original product, alternative, savings, purchase link

**Savings Parser (line 1276):**
- Updated delimiter detection for extracting savings amount
- Ensures minimum $5 filter continues to work
- Compatible with new dash-based format

#### Live Example Generated:
```
"Echo Dot 5th Gen ($49.99) - Alternative: Google Nest Mini ($29.99) - Save: $20.00 - Buy: https://store.google.com/product/google_nest_mini"
```

**Backend Process ID**: bd2be2f
**Status**: Successfully generating notifications with new format ✅
**Output Log**: `/tmp/claude/-home-sholto-Documents-Hobbies-hackathons/tasks/bd2be2f.output`

---

**Session Status**: Both servers running successfully ✅
**Backend**: Actively checking purchases for cheaper alternatives every 30s (Process: bd2be2f)
**Transactions**: All 60 mock transactions displayed with correct arrows
**Links**: All recommendations include clickable purchase links
**Filter**: Only showing alternatives with >= $5 savings
**Format**: Clean, readable format without "You bought:" prefix or vertical pipes
**Document Created**: January 31, 2026, 20:33 UTC
**Last Updated**: January 31, 2026, 22:35 UTC
