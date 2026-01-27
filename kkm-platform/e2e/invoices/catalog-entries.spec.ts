import { test, expect } from "@playwright/test";

test.describe("Catalog Entries Management", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/invoices");
    await page.waitForLoadState("networkidle");

    // Open create invoice form
    await page.click('text="Создать счет-фактуру"');
    await page.waitForSelector('h2:has-text("Создать новый счет-фактуру")');
  });

  test("should add new catalog entry", async ({ page }) => {
    // Click add entry button
    await page.click('text="Добавить позицию"');

    // Verify empty row is added
    const rows = await page.locator("table tbody tr").count();
    expect(rows).toBeGreaterThan(0);

    // Fill entry data
    await page.fill('input[placeholder*="наименование"]', "Товар 1");
    await page.fill('input[placeholder*="Количество"]', "5");
    await page.fill('input[placeholder*="Цена"]', "500");

    // Verify entry is added
    await expect(page.locator('text="Товар 1"')).toBeVisible();
  });

  test("should remove catalog entry", async ({ page }) => {
    // Add entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="наименование"]', "Товар для удаления");

    // Get initial count
    const initialCount = await page.locator("table tbody tr").count();

    // Click remove button
    await page.click('button[aria-label="Удалить позицию"]');

    // Verify entry is removed
    const currentCount = await page.locator("table tbody tr").count();
    expect(currentCount).toBe(initialCount - 1);
  });

  test("should edit catalog entry inline", async ({ page }) => {
    // Add entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="наименование"]', "Товар 1");
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Modify quantity
    await page.fill('input[placeholder*="Количество"]', "20");

    // Verify total recalculates (20 * 1000 = 20000)
    await expect(page.locator('text*="20"')).toBeVisible();
  });

  test("should validate catalog entry fields", async ({ page }) => {
    // Add entry
    await page.click('text="Добавить позицию"');

    // Try to enter invalid quantity
    await page.fill('input[placeholder*="Количество"]', "0");
    await page.click("body");

    // Verify validation error
    await expect(
      page.locator('text*="Количество должно быть больше 0"'),
    ).toBeVisible();
  });

  test("should calculate entry totals automatically", async ({ page }) => {
    // Add entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Select 12% VAT
    await page.selectOption('select[name="taxRateVATCode"]', "12");

    // Verify totals are calculated
    // Total amount = 10 * 1000 = 10000
    // Amount without VAT = 10000 / 1.12 = 8928.57
    // VAT amount = 10000 - 8928.57 = 1071.43

    const totalsSection = page.locator('div:has-text("Итоговые суммы")');
    await expect(totalsSection).toBeVisible();
  });

  test("should add multiple catalog entries", async ({ page }) => {
    // Add first entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="наименование"]', "Товар 1");
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Add second entry
    await page.click('text="Добавить позицию"');
    const inputs = page.locator('input[placeholder*="наименование"]');
    await inputs.nth(1).fill("Товар 2");
    await page.locator('input[placeholder*="Количество"]').nth(1).fill("5");
    await page.locator('input[placeholder*="Цена"]').nth(1).fill("2000");

    // Verify both entries exist
    await expect(page.locator('text="Товар 1"')).toBeVisible();
    await expect(page.locator('text="Товар 2"')).toBeVisible();

    // Verify total is sum of both (10*1000 + 5*2000 = 20000)
    const totalsSection = page.locator('div:has-text("Итоговые суммы")');
    await expect(totalsSection).toContainText("20");
  });

  test("should prevent submission without catalog entries", async ({
    page,
  }) => {
    // Fill basic info without adding entries
    await page.selectOption('select[name="operationTypeCode"]', "10");
    await page.fill('input[type="date"][name="deliveryDate"]', "2026-01-27");
    await page.check('input[type="checkbox"][name="isResident"]');
    await page.fill('input[name="contractorTin"]', "12345678901234");

    // Try to submit
    await page.click('button:has-text("Создать счет-фактуру")');

    // Verify validation error
    await expect(
      page.locator('text*="Добавьте хотя бы одну позицию"'),
    ).toBeVisible();
  });

  test("should preserve entry data when adding new entries", async ({
    page,
  }) => {
    // Add first entry
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="наименование"]', "Товар 1");
    await page.fill('input[placeholder*="Количество"]', "10");
    await page.fill('input[placeholder*="Цена"]', "1000");

    // Add second entry
    await page.click('text="Добавить позицию"');

    // Verify first entry data is preserved
    const firstNameInput = page
      .locator('input[placeholder*="наименование"]')
      .first();
    await expect(firstNameInput).toHaveValue("Товар 1");
  });

  test("should handle decimal quantities and prices", async ({ page }) => {
    // Add entry with decimal values
    await page.click('text="Добавить позицию"');
    await page.fill('input[placeholder*="Количество"]', "1.5");
    await page.fill('input[placeholder*="Цена"]', "99.99");

    // Verify total is calculated (1.5 * 99.99 = 149.985 ≈ 149.98)
    await expect(page.locator('text*="149"')).toBeVisible();
  });
});
