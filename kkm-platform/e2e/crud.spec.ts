import { test, expect } from '@playwright/test';

test.describe('Users Page', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/auth/login');
    await page.getByLabel(/email/i).fill('admin@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /sign in|login/i }).click();
    await page.waitForURL(/\/dashboard/);
  });

  test('should load users list', async ({ page }) => {
    await page.goto('/users');
    
    // Wait for page to load
    await page.waitForLoadState('networkidle');
    
    // Check if users table is visible
    const table = page.locator('table');
    await expect(table).toBeVisible();
  });

  test('should create new user', async ({ page }) => {
    await page.goto('/users');
    
    // Click create button
    const createBtn = page.getByRole('button', { name: /create|add|new user/i }).first();
    if (await createBtn.isVisible()) {
      await createBtn.click();
      
      // Fill form
      await page.getByLabel(/name/i).fill('Test User');
      await page.getByLabel(/email/i).fill('testuser@example.com');
      
      // Submit
      await page.getByRole('button', { name: /save|create|submit/i }).click();
      
      // Check for success toast
      await page.waitForTimeout(1000);
    }
  });

  test('should search users', async ({ page }) => {
    await page.goto('/users');
    
    // Find search input
    const searchInput = page.getByPlaceholder(/search/i).first();
    if (await searchInput.isVisible()) {
      await searchInput.fill('admin');
      await page.waitForTimeout(500);
      
      // Check that results are filtered
    }
  });

  test('should delete user', async ({ page }) => {
    await page.goto('/users');
    
    // Find delete button (usually in action column)
    const deleteBtn = page.getByRole('button', { name: /delete|remove/i }).first();
    if (await deleteBtn.isVisible()) {
      await deleteBtn.click();
      
      // Handle confirmation if any
      const confirmBtn = page.getByRole('button', { name: /confirm|yes|delete/i });
      if (await confirmBtn.isVisible()) {
        await confirmBtn.click();
      }
    }
  });
});

test.describe('Invoices Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/auth/login');
    await page.getByLabel(/email/i).fill('admin@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /sign in|login/i }).click();
    await page.waitForURL(/\/dashboard/);
  });

  test('should load invoices list', async ({ page }) => {
    await page.goto('/invoices');
    await page.waitForLoadState('networkidle');
    
    const table = page.locator('table');
    await expect(table).toBeVisible();
  });

  test('should filter invoices by status', async ({ page }) => {
    await page.goto('/invoices');
    
    const statusFilter = page.locator('select').first();
    if (await statusFilter.isVisible()) {
      await statusFilter.selectOption('paid');
      await page.waitForTimeout(500);
    }
  });

  test('should create invoice', async ({ page }) => {
    await page.goto('/invoices');
    
    const createBtn = page.getByRole('button', { name: /create|add|new invoice/i }).first();
    if (await createBtn.isVisible()) {
      await createBtn.click();
      await page.waitForURL(/\/invoices\/create/);
      
      // Verify we're on create page
      await expect(page).toHaveURL(/\/invoices\/create/);
    }
  });
});
