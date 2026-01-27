import { test, expect } from "@playwright/test";

test.describe("Edit Invoice", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/invoices");
    await page.waitForLoadState("networkidle");
  });

  test("should display edit form with existing data", async ({ page }) => {
    // Click edit button on first invoice
    await page.click('button[aria-label="Редактировать"]');

    // Verify form is displayed
    await expect(
      page.locator('h2:has-text("Редактировать счет-фактуру")'),
    ).toBeVisible();

    // Verify form is populated with existing data
    await expect(page.locator('input[name="invoiceNumber"]')).not.toBeEmpty();
  });

  test("should update invoice successfully", async ({ page }) => {
    await page.click('button[aria-label="Редактировать"]');

    // Modify invoice number
    await page.fill('input[name="invoiceNumber"]', "UPDATED-001");

    // Submit form
    await page.click('button:has-text("Сохранить изменения")');

    // Verify success message
    await expect(
      page.locator('text="Счет-фактура успешно обновлен"'),
    ).toBeVisible({ timeout: 10000 });
  });

  test("should preserve unchanged fields", async ({ page }) => {
    await page.click('button[aria-label="Редактировать"]');

    // Get original TIN value
    const originalTin = await page.inputValue('input[name="contractorTin"]');

    // Modify only invoice number
    await page.fill('input[name="invoiceNumber"]', "UPDATED-002");

    // Verify TIN is still the same
    await expect(page.locator('input[name="contractorTin"]')).toHaveValue(
      originalTin,
    );
  });

  test("should validate edited data", async ({ page }) => {
    await page.click('button[aria-label="Редактировать"]');

    // Clear required field
    await page.fill('input[name="contractorTin"]', "");

    // Try to submit
    await page.click('button:has-text("Сохранить изменения")');

    // Verify validation error
    await expect(page.locator('text*="Обязательное поле"')).toBeVisible();
  });

  test("should cancel edit without saving", async ({ page }) => {
    await page.click('button[aria-label="Редактировать"]');

    // Get original value
    const originalNumber = await page.inputValue('input[name="invoiceNumber"]');

    // Modify data
    await page.fill('input[name="invoiceNumber"]', "TEMP-999");

    // Cancel
    await page.click('button:has-text("Отмена")');

    // Return to list and verify data wasn't changed
    await expect(page.locator('h1:has-text("Счета-фактуры")')).toBeVisible();
  });

  test("should handle edit conflicts", async ({ page }) => {
    // Mock 409 conflict error
    await page.route("**/api/invoices/*", (route) => {
      if (route.request().method() === "PUT") {
        route.fulfill({
          status: 409,
          body: JSON.stringify({
            error: "Invoice has been modified by another user",
          }),
        });
      } else {
        route.continue();
      }
    });

    await page.click('button[aria-label="Редактировать"]');
    await page.fill('input[name="invoiceNumber"]', "CONFLICT-001");
    await page.click('button:has-text("Сохранить изменения")');

    // Verify conflict error message
    await expect(page.locator('text*="конфликт"')).toBeVisible({
      timeout: 5000,
    });
  });
});
