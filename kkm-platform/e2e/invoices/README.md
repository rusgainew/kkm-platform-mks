# E2E Tests for Invoice Feature

Comprehensive end-to-end tests for the Invoice Management system using Playwright.

## Test Structure

### 1. Create Invoice Tests (`create-invoice.spec.ts`)

- ✅ Display create invoice form
- ✅ Create invoice with valid data
- ✅ Validate required fields
- ✅ Validate TIN format
- ✅ Calculate totals automatically
- ✅ Clear form on cancel
- ✅ Handle API errors gracefully

**Total: 7 tests**

### 2. Edit Invoice Tests (`edit-invoice.spec.ts`)

- ✅ Display edit form with existing data
- ✅ Update invoice successfully
- ✅ Preserve unchanged fields
- ✅ Validate edited data
- ✅ Cancel edit without saving
- ✅ Handle edit conflicts

**Total: 6 tests**

### 3. Delete Invoice Tests (`delete-invoice.spec.ts`)

- ✅ Display delete confirmation dialog
- ✅ Delete invoice successfully
- ✅ Cancel deletion
- ✅ Handle deletion errors
- ✅ Disable delete for certain invoice statuses

**Total: 5 tests**

### 4. Catalog Entries Tests (`catalog-entries.spec.ts`)

- ✅ Add new catalog entry
- ✅ Remove catalog entry
- ✅ Edit catalog entry inline
- ✅ Validate catalog entry fields
- ✅ Calculate entry totals automatically
- ✅ Add multiple catalog entries
- ✅ Prevent submission without catalog entries
- ✅ Preserve entry data when adding new entries
- ✅ Handle decimal quantities and prices

**Total: 9 tests**

## Running Tests

### Run all invoice E2E tests:

```bash
pnpm exec playwright test e2e/invoices/
```

### Run specific test file:

```bash
pnpm exec playwright test e2e/invoices/create-invoice.spec.ts
```

### Run in headed mode (with browser UI):

```bash
pnpm exec playwright test e2e/invoices/ --headed
```

### Run with debugging:

```bash
pnpm exec playwright test e2e/invoices/ --debug
```

### Run single test:

```bash
pnpm exec playwright test e2e/invoices/ --grep "should create invoice with valid data"
```

## Test Coverage

**Total E2E Tests: 27**

- Create Invoice: 7 tests
- Edit Invoice: 6 tests
- Delete Invoice: 5 tests
- Catalog Entries: 9 tests

## Test Features

1. **Form Validation**
   - Required field validation
   - TIN format validation
   - Catalog entry validation

2. **CRUD Operations**
   - Create new invoices
   - Edit existing invoices
   - Delete invoices
   - Manage catalog entries

3. **Error Handling**
   - API errors
   - Network failures
   - Conflict resolution

4. **Calculations**
   - Automatic VAT calculation
   - Total amount calculation
   - Entry-level calculations

5. **User Interactions**
   - Form submission
   - Cancel actions
   - Confirmation dialogs
   - Inline editing

## Prerequisites

- Playwright browsers installed
- Next.js dev server running
- API server available at localhost:8080

## Browser Support

Tests run on:

- Chromium (primary)
- Firefox (optional)
- WebKit (optional)

To install browsers:

```bash
npx playwright install chromium
```

## Notes

- Tests use real browser automation
- Tests interact with actual UI components
- API calls can be mocked for isolated testing
- Tests use Russian locale for UI elements
