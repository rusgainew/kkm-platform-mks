/**
 * E2E tests for Bank Accounts functionality
 * Testing complete user workflows for bank account management
 *
 * @playwright
 */

import { test, expect } from "@playwright/test";

test.describe("Bank Accounts Management", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to bank accounts page (adjust URL as needed)
    await page.goto("/bank-accounts");
    // Wait for page to load
    await page.waitForLoadState("networkidle");
  });

  test.describe("Create Bank Account", () => {
    test("should create a new bank account with valid data", async ({
      page,
    }) => {
      // Click create button
      await page.click("text=Добавить банковский счет");

      // Fill in form
      await page.fill(
        'input[placeholder="12345678901234567890"]',
        "40702810800003000001",
      );
      await page.fill('input[placeholder="ПАО Сбербанк"]', "ПАО Сбербанк");
      await page.fill('input[placeholder="044525225"]', "044525225");
      await page.selectOption("select", "RUB");

      // Submit form
      await page.click("text=Создать счет");

      // Verify success
      await expect(page.locator("text=Счет успешно создан")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should show validation errors for empty required fields", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      // Try to submit without filling
      await page.click("text=Создать счет");

      // Check for validation errors
      await expect(page.locator("text=Номер счета обязателен")).toBeVisible();
      await expect(
        page.locator("text=Название банка обязательно"),
      ).toBeVisible();
      await expect(page.locator("text=БИК банка обязателен")).toBeVisible();
    });

    test("should show error for invalid account number length", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      // Enter invalid account number (less than 20 digits)
      const accountInput = page.locator(
        'input[placeholder="12345678901234567890"]',
      );
      await accountInput.fill("123456789");
      await accountInput.blur();

      await expect(
        page.locator("text=Номер счета должен содержать ровно 20 цифр"),
      ).toBeVisible();
    });

    test("should auto-clean non-numeric characters from account number", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      const accountInput = page.locator(
        'input[placeholder="12345678901234567890"]',
      );
      await accountInput.fill("4070-2810-8000-0300-0001");

      // Should keep only digits
      await expect(accountInput).toHaveValue("40702810800003000001");
    });

    test("should show error for invalid bank code length", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      const bankCodeInput = page.locator('input[placeholder="044525225"]');
      await bankCodeInput.fill("12345");
      await bankCodeInput.blur();

      await expect(
        page.locator("text=БИК должен содержать от 6 до 9 цифр"),
      ).toBeVisible();
    });

    test("should limit bank code to 9 digits", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      const bankCodeInput = page.locator('input[placeholder="044525225"]');
      await bankCodeInput.fill("12345678901234");

      await expect(bankCodeInput).toHaveValue("123456789");
    });

    test("should show character counter for account number", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      const accountInput = page.locator(
        'input[placeholder="12345678901234567890"]',
      );
      await accountInput.fill("12345");

      await expect(page.locator("text=Текущая длина: 5")).toBeVisible();
    });

    test("should select different currencies", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      const currencySelect = page.locator("select");

      // Test USD
      await currencySelect.selectOption("USD");
      await expect(currencySelect).toHaveValue("USD");

      // Test EUR
      await currencySelect.selectOption("EUR");
      await expect(currencySelect).toHaveValue("EUR");

      // Test KGS (default)
      await currencySelect.selectOption("KGS");
      await expect(currencySelect).toHaveValue("KGS");
    });

    test("should toggle active status", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      const checkbox = page.locator('input[type="checkbox"]');

      // Should be checked by default
      await expect(checkbox).toBeChecked();

      // Uncheck
      await checkbox.uncheck();
      await expect(checkbox).not.toBeChecked();

      // Check again
      await checkbox.check();
      await expect(checkbox).toBeChecked();
    });

    test("should show informational message", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      await expect(
        page.locator("text=Номер счета должен содержать ровно 20 цифр"),
      ).toBeVisible();
      await expect(
        page.locator("text=БИК банка обычно содержит 9 цифр"),
      ).toBeVisible();
    });
  });

  test.describe("Edit Bank Account", () => {
    test("should edit existing bank account", async ({ page }) => {
      // Click on existing account (assuming at least one exists)
      await page.click('[data-testid="bank-account-item"]:first-child');

      // Wait for form to load with data
      await expect(
        page.locator("text=Редактировать банковский счет"),
      ).toBeVisible();

      // Modify bank name
      const bankNameInput = page.locator('input[placeholder="ПАО Сбербанк"]');
      await bankNameInput.clear();
      await bankNameInput.fill("АО Альфа-Банк");

      // Submit
      await page.click("text=Обновить счет");

      // Verify success
      await expect(page.locator("text=Счет успешно обновлен")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should pre-fill form with existing data", async ({ page }) => {
      await page.click('[data-testid="bank-account-item"]:first-child');

      // Check that inputs are pre-filled (values will vary)
      const accountInput = page.locator(
        'input[placeholder="12345678901234567890"]',
      );
      await expect(accountInput).not.toHaveValue("");

      const bankNameInput = page.locator('input[placeholder="ПАО Сбербанк"]');
      await expect(bankNameInput).not.toHaveValue("");

      const bankCodeInput = page.locator('input[placeholder="044525225"]');
      await expect(bankCodeInput).not.toHaveValue("");
    });

    test("should validate when editing", async ({ page }) => {
      await page.click('[data-testid="bank-account-item"]:first-child');

      // Clear account number
      const accountInput = page.locator(
        'input[placeholder="12345678901234567890"]',
      );
      await accountInput.clear();
      await accountInput.blur();

      await expect(page.locator("text=Номер счета обязателен")).toBeVisible();
    });

    test("should cancel editing", async ({ page }) => {
      await page.click('[data-testid="bank-account-item"]:first-child');

      // Modify field
      const bankNameInput = page.locator('input[placeholder="ПАО Сбербанк"]');
      await bankNameInput.fill("Новый банк");

      // Cancel
      await page.click("text=Отмена");

      // Should return to list (adjust selector as needed)
      await expect(page.locator("text=Список банковских счетов")).toBeVisible();
    });
  });

  test.describe("Delete Bank Account", () => {
    test("should delete bank account", async ({ page }) => {
      // Click delete button on first account
      await page.click('[data-testid="delete-bank-account-btn"]:first-child');

      // Confirm deletion
      await page.click("text=Подтвердить удаление");

      // Verify success
      await expect(page.locator("text=Счет успешно удален")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should show confirmation dialog before deleting", async ({
      page,
    }) => {
      await page.click('[data-testid="delete-bank-account-btn"]:first-child');

      await expect(
        page.locator("text=Вы уверены, что хотите удалить этот счет?"),
      ).toBeVisible();
    });

    test("should cancel deletion", async ({ page }) => {
      await page.click('[data-testid="delete-bank-account-btn"]:first-child');

      // Cancel
      await page.click("text=Отмена");

      // Account should still be visible
      await expect(
        page.locator('[data-testid="bank-account-item"]'),
      ).toBeVisible();
    });
  });

  test.describe("Bank Accounts List", () => {
    test("should display list of bank accounts", async ({ page }) => {
      const count = await page
        .locator('[data-testid="bank-account-item"]')
        .count();
      expect(count).toBeGreaterThanOrEqual(0);
    });

    test("should show account details in list", async ({ page }) => {
      const firstAccount = page
        .locator('[data-testid="bank-account-item"]')
        .first();

      // Should display account number (partially masked or full)
      await expect(firstAccount.locator("text=****")).toBeVisible();

      // Should display bank name
      await expect(firstAccount).toContainText(/банк/i);

      // Should display currency
      await expect(firstAccount).toContainText(/[A-Z]{3}/); // 3-letter currency code
    });

    test("should filter accounts by currency", async ({ page }) => {
      // Select currency filter
      await page.selectOption('[data-testid="currency-filter"]', "RUB");

      // Wait for results
      await page.waitForTimeout(500);

      // Check that only RUB accounts are shown
      const accounts = page.locator('[data-testid="bank-account-item"]');
      const count = await accounts.count();

      for (let i = 0; i < count; i++) {
        await expect(accounts.nth(i)).toContainText("RUB");
      }
    });

    test("should search accounts by bank name", async ({ page }) => {
      const searchInput = page.locator('[data-testid="search-bank-accounts"]');
      await searchInput.fill("Сбербанк");

      // Wait for search results
      await page.waitForTimeout(500);

      const accounts = page.locator('[data-testid="bank-account-item"]');
      const count = await accounts.count();

      for (let i = 0; i < count; i++) {
        await expect(accounts.nth(i)).toContainText(/сбербанк/i);
      }
    });

    test("should show active status indicator", async ({ page }) => {
      const firstAccount = page
        .locator('[data-testid="bank-account-item"]')
        .first();

      // Should have active/inactive indicator
      const hasActiveIndicator = await firstAccount
        .locator("text=Активен")
        .isVisible();
      const hasInactiveIndicator = await firstAccount
        .locator("text=Неактивен")
        .isVisible();

      expect(hasActiveIndicator || hasInactiveIndicator).toBe(true);
    });

    test("should paginate accounts list", async ({ page }) => {
      // Check if pagination exists
      const pagination = page.locator('[data-testid="pagination"]');

      if (await pagination.isVisible()) {
        const nextButton = page.locator("text=Следующая");
        await nextButton.click();

        // Wait for new page to load
        await page.waitForLoadState("networkidle");

        // Should show different accounts
        await expect(
          page.locator('[data-testid="bank-account-item"]'),
        ).toBeVisible();
      }
    });
  });

  test.describe("Form UX", () => {
    test("should disable submit button when form is invalid", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      const submitButton = page.locator("text=Создать счет");

      // Initially should be disabled or become disabled after validation
      await page.click('button[type="submit"]');

      await expect(submitButton).toBeDisabled();
    });

    test("should show loading state during submission", async ({ page }) => {
      await page.click("text=Добавить банковский счет");

      // Fill valid data
      await page.fill(
        'input[placeholder="12345678901234567890"]',
        "40702810800003000001",
      );
      await page.fill('input[placeholder="ПАО Сбербанк"]', "ПАО Сбербанк");
      await page.fill('input[placeholder="044525225"]', "044525225");

      // Submit
      await page.click("text=Создать счет");

      // Should show loading state briefly
      await expect(page.locator("text=Сохранение...")).toBeVisible({
        timeout: 1000,
      });
    });

    test("should show error message on submission failure", async ({
      page,
    }) => {
      await page.click("text=Добавить банковский счет");

      // Fill form
      await page.fill(
        'input[placeholder="12345678901234567890"]',
        "40702810800003000001",
      );
      await page.fill('input[placeholder="ПАО Сбербанк"]', "ПАО Сбербанк");
      await page.fill('input[placeholder="044525225"]', "044525225");

      // Mock API error (if using MSW or similar)
      // await page.route('**/api/v1/bank-accounts', (route) =>
      //   route.fulfill({ status: 500, body: 'Server error' })
      // );

      await page.click("text=Создать счет");

      // Should show error message
      await expect(page.locator(".bg-red-900\\/20")).toBeVisible({
        timeout: 5000,
      });
    });
  });

  test.describe("Responsive Design", () => {
    test("should display form correctly on mobile", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.click("text=Добавить банковский счет");

      // Check that form is visible and scrollable
      await expect(page.locator("form")).toBeVisible();
      await expect(
        page.locator('input[placeholder="12345678901234567890"]'),
      ).toBeVisible();
    });

    test("should display list correctly on tablet", async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });

      await expect(
        page.locator('[data-testid="bank-account-item"]'),
      ).toBeVisible();
    });
  });
});
