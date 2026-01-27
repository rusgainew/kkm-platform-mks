import { test, expect } from "@playwright/test";

test.describe("Delete Invoice", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/invoices");
    await page.waitForLoadState("networkidle");
  });

  test("should display delete confirmation dialog", async ({ page }) => {
    // Click delete button
    await page.click('button[aria-label="Удалить"]');

    // Verify confirmation dialog
    await expect(page.locator('text*="Вы уверены"')).toBeVisible();
    await expect(page.locator('button:has-text("Удалить")')).toBeVisible();
    await expect(page.locator('button:has-text("Отмена")')).toBeVisible();
  });

  test("should delete invoice successfully", async ({ page }) => {
    // Get invoice number before deletion
    const invoiceNumber = await page
      .locator("table tbody tr:first-child td:nth-child(2)")
      .textContent();

    // Click delete
    await page.click('button[aria-label="Удалить"]');

    // Confirm deletion
    await page.click('button:has-text("Удалить")');

    // Verify success message
    await expect(
      page.locator('text="Счет-фактура успешно удален"'),
    ).toBeVisible({ timeout: 10000 });

    // Verify invoice is removed from list
    await expect(page.locator(`text="${invoiceNumber}"`)).not.toBeVisible();
  });

  test("should cancel deletion", async ({ page }) => {
    // Get initial row count
    const initialCount = await page.locator("table tbody tr").count();

    // Click delete
    await page.click('button[aria-label="Удалить"]');

    // Cancel deletion
    await page.click('button:has-text("Отмена")');

    // Verify dialog is closed
    await expect(page.locator('text*="Вы уверены"')).not.toBeVisible();

    // Verify row count unchanged
    const currentCount = await page.locator("table tbody tr").count();
    expect(currentCount).toBe(initialCount);
  });

  test("should handle deletion errors", async ({ page }) => {
    // Mock API error
    await page.route("**/api/invoices/*", (route) => {
      if (route.request().method() === "DELETE") {
        route.fulfill({
          status: 500,
          body: JSON.stringify({ error: "Cannot delete invoice" }),
        });
      } else {
        route.continue();
      }
    });

    await page.click('button[aria-label="Удалить"]');
    await page.click('button:has-text("Удалить")');

    // Verify error message
    await expect(page.locator('text*="Ошибка"')).toBeVisible({ timeout: 5000 });
  });

  test("should disable delete for certain invoice statuses", async ({
    page,
  }) => {
    // Mock invoice with non-deletable status
    await page.route("**/api/invoices", (route) => {
      route.fulfill({
        status: 200,
        body: JSON.stringify({
          invoices: [
            {
              id: 1,
              invoiceNumber: "INV-001",
              status: "approved",
              canDelete: false,
            },
          ],
        }),
      });
    });

    await page.reload();

    // Verify delete button is disabled
    const deleteButton = page.locator('button[aria-label="Удалить"]');
    await expect(deleteButton).toBeDisabled();
  });
});
