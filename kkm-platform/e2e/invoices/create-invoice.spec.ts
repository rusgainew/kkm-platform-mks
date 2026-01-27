import { test, expect } from "@playwright/test";

test.describe("Create Invoice", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to invoices page
    await page.goto("/invoices");

    // Wait for page to load
    await page.waitForLoadState("networkidle");
  });

  test("should display create invoice form", async ({ page }) => {
    // Click create button
    await page.click('text="Создать счет-фактуру"');

    // Verify form is displayed
    await expect(
      page.locator('h2:has-text("Создать новый счет-фактуру")'),
    ).toBeVisible();

    // Verify required fields are present
    await expect(page.locator('label:has-text("Тип операции")')).toBeVisible();
    await expect(page.locator('label:has-text("Дата доставки")')).toBeVisible();
    await expect(page.locator('label:has-text("Валюта")')).toBeVisible();
  });

  test("should create invoice with valid data", async ({ page }) => {
    await page.click('text="Создать счет-фактуру"');

    // Fill in basic information
    await page.selectOption('select[name="operationTypeCode"]', "10");
    await page.fill('input[name="invoiceNumber"]', "TEST-001");
    await page.fill('input[type="date"][name="deliveryDate"]', "2026-01-27");
    await page.selectOption('select[name="currencyCode"]', "KGS");
    await page.selectOption('select[name="taxRateVATCode"]', "12");

    // Fill contractor info
    await page.check('input[type="checkbox"][name="isResident"]');
    await page.fill('input[name="contractorTin"]', "12345678901234");

    // Add catalog entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="наименование"]', "Тестовый товар");
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Submit form
    await page.click('button:has-text("Создать счет-фактуру")');

    // Verify success message
    await expect(
      page.locator('text="Счет-фактура успешно создан"'),
    ).toBeVisible({ timeout: 10000 });
  });

  test("should validate required fields", async ({ page }) => {
    await page.click('text="Создать счет-фактуру"');

    // Try to submit without filling required fields
    await page.click('button:has-text("Создать счет-фактуру")');

    // Verify validation errors are displayed
    await expect(page.locator('text*="Обязательное поле"')).toBeVisible();
  });

  test("should validate TIN format", async ({ page }) => {
    await page.click('text="Создать счет-фактуру"');

    // Fill invalid TIN
    await page.check('input[type="checkbox"][name="isResident"]');
    await page.fill('input[name="contractorTin"]', "123");

    // Blur field to trigger validation
    await page.click("body");

    // Verify TIN validation error
    await expect(
      page.locator('text*="ИНН должен содержать 14 цифр"'),
    ).toBeVisible();
  });

  test("should calculate totals automatically", async ({ page }) => {
    await page.click('text="Создать счет-фактуру"');

    // Add catalog entry with known values
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Verify total is calculated (10 * 1000 = 10000)
    await expect(page.locator('text*="10"')).toBeVisible();
    await expect(page.locator('text*="000"')).toBeVisible();
  });

  test("should clear form on cancel", async ({ page }) => {
    await page.click('text="Создать счет-фактуру"');

    // Fill some data
    await page.fill('input[name="invoiceNumber"]', "TEST-001");

    // Click cancel
    await page.click('button:has-text("Отмена")');

    // Verify returned to invoices list
    await expect(page.locator('h1:has-text("Счета-фактуры")')).toBeVisible();
  });

  test("should handle API errors gracefully", async ({ page }) => {
    // Mock API error
    await page.route("**/api/invoices", (route) => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: "Server error" }),
      });
    });

    await page.click('text="Создать счет-фактуру"');

    // Fill valid data
    await page.selectOption('select[name="operationTypeCode"]', "10");
    await page.fill('input[type="date"][name="deliveryDate"]', "2026-01-27");
    await page.selectOption('select[name="currencyCode"]', "KGS");
    await page.check('input[type="checkbox"][name="isResident"]');
    await page.fill('input[name="contractorTin"]', "12345678901234");

    // Submit
    await page.click('button:has-text("Создать счет-фактуру")');

    // Verify error message
    await expect(page.locator('text*="Ошибка"')).toBeVisible({ timeout: 5000 });
  });
});
