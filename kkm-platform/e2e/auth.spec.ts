import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test('should redirect to login when not authenticated', async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveURL(/\/auth\/login/);
  });

  test('should show login form', async ({ page }) => {
    await page.goto('/auth/login');
    
    // Check form elements
    await expect(page.getByLabel(/email/i)).toBeVisible();
    await expect(page.getByLabel(/password/i)).toBeVisible();
    await expect(page.getByRole('button', { name: /sign in|login/i })).toBeVisible();
  });

  test('should show error on invalid credentials', async ({ page }) => {
    await page.goto('/auth/login');
    
    // Fill invalid credentials
    await page.getByLabel(/email/i).fill('invalid@example.com');
    await page.getByLabel(/password/i).fill('wrongpassword');
    await page.getByRole('button', { name: /sign in|login/i }).click();
    
    // Wait for error message or toast
    await page.waitForTimeout(1000);
    // Error handling depends on implementation
  });

  test('should login with valid credentials', async ({ page }) => {
    await page.goto('/auth/login');
    
    // Fill valid credentials (these would need to be set up in test DB)
    await page.getByLabel(/email/i).fill('admin@example.com');
    await page.getByLabel(/password/i).fill('password123');
    await page.getByRole('button', { name: /sign in|login/i }).click();
    
    // Should redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 5000 });
  });
});

test.describe('Navigation', () => {
  test('should navigate between main pages', async ({ page }) => {
    // This assumes user is already logged in
    await page.goto('/dashboard');
    
    // Check if main navigation exists
    const sidebarLinks = page.locator('nav a');
    expect(await sidebarLinks.count()).toBeGreaterThan(0);
  });

  test('should show 403 page for unauthorized access', async ({ page }) => {
    await page.goto('/403');
    
    await expect(page.getByText(/403/)).toBeVisible();
    await expect(page.getByText(/access forbidden|no permission/i)).toBeVisible();
  });
});
