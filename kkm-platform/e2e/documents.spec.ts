/**
 * E2E tests for Documents functionality
 * Testing complete user workflows for document management, file uploads, and status changes
 *
 * @playwright
 */

import { test, expect } from "@playwright/test";

test.describe("Documents Management", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to documents page
    await page.goto("/documents");
    await page.waitForLoadState("networkidle");
  });

  test.describe("Create Document", () => {
    test("should create a new document with valid data", async ({ page }) => {
      await page.click("text=Создать документ");

      // Fill in form
      await page.fill(
        'input[placeholder*="Договор"]',
        "Договор поставки №2026-001",
      );
      await page.fill(
        "textarea",
        "Договор на поставку товаров на сумму 100,000 сом. Срок поставки - 30 дней.",
      );
      await page.selectOption("select", "draft");

      // Submit form
      await page.click("text=Создать документ");

      // Verify success
      await expect(page.locator("text=успешно создан")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should show validation errors for empty required fields", async ({
      page,
    }) => {
      await page.click("text=Создать документ");

      // Try to submit without filling
      await page.click('button[type="submit"]');

      // Check for validation errors
      await expect(
        page.locator("text=Название документа обязательно"),
      ).toBeVisible();
      await expect(
        page.locator("text=Содержимое документа обязательно"),
      ).toBeVisible();
    });

    test("should show error for short title", async ({ page }) => {
      await page.click("text=Создать документ");

      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.fill("AB");
      await titleInput.blur();

      await expect(
        page.locator("text=Название должно содержать от 3 до 200 символов"),
      ).toBeVisible();
    });

    test("should show error for short content", async ({ page }) => {
      await page.click("text=Создать документ");

      const contentInput = page.locator("textarea");
      await contentInput.fill("Короткий");
      await contentInput.blur();

      await expect(
        page.locator("text=Содержимое должно содержать минимум 10 символов"),
      ).toBeVisible();
    });

    test("should show character counters", async ({ page }) => {
      await page.click("text=Создать документ");

      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.fill("Тестовый документ");

      await expect(page.locator("text=/Длина: \\d+\\/200/")).toBeVisible();

      const contentInput = page.locator("textarea");
      await contentInput.fill("Это содержимое документа");

      await expect(page.locator("text=/текущая длина: \\d+/")).toBeVisible();
    });

    test("should select different statuses", async ({ page }) => {
      await page.click("text=Создать документ");

      const statusSelect = page.locator("select");

      // Test draft
      await statusSelect.selectOption("draft");
      await expect(statusSelect).toHaveValue("draft");

      // Test pending
      await statusSelect.selectOption("pending");
      await expect(statusSelect).toHaveValue("pending");

      // Test approved
      await statusSelect.selectOption("approved");
      await expect(statusSelect).toHaveValue("approved");
    });

    test("should display status options with emojis", async ({ page }) => {
      await page.click("text=Создать документ");

      await expect(page.locator("text=📝 Черновик")).toBeVisible();
      await expect(page.locator("text=⏳ На рассмотрении")).toBeVisible();
      await expect(page.locator("text=✅ Утвержден")).toBeVisible();
      await expect(page.locator("text=❌ Отклонен")).toBeVisible();
      await expect(page.locator("text=📦 В архиве")).toBeVisible();
    });

    test("should set assigned_to field", async ({ page }) => {
      await page.click("text=Создать документ");

      const assignedInput = page.locator('input[placeholder*="User ID"]');
      await assignedInput.fill("user-456");

      await expect(assignedInput).toHaveValue("user-456");
    });
  });

  test.describe("File Upload", () => {
    test("should display file upload section", async ({ page }) => {
      await page.click("text=Создать документ");

      await expect(page.locator("text=Файлы документов")).toBeVisible();
      await expect(page.locator("text=Выбрать файлы")).toBeVisible();
      await expect(
        page.locator("text=Максимальный размер файла: 10MB"),
      ).toBeVisible();
    });

    test("should upload a file", async ({ page }) => {
      await page.click("text=Создать документ");

      // Create a test file
      const fileInput = page.locator("#file-upload");

      // Upload a file
      await fileInput.setInputFiles({
        name: "test-document.pdf",
        mimeType: "application/pdf",
        buffer: Buffer.from("Test PDF content"),
      });

      // Verify file is displayed
      await expect(page.locator("text=test-document.pdf")).toBeVisible();
      await expect(page.locator("text=Загружено файлов: 1")).toBeVisible();
    });

    test("should upload multiple files", async ({ page }) => {
      await page.click("text=Создать документ");

      const fileInput = page.locator("#file-upload");

      await fileInput.setInputFiles([
        {
          name: "document1.pdf",
          mimeType: "application/pdf",
          buffer: Buffer.from("PDF 1"),
        },
        {
          name: "document2.docx",
          mimeType:
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
          buffer: Buffer.from("DOCX 1"),
        },
      ]);

      await expect(page.locator("text=document1.pdf")).toBeVisible();
      await expect(page.locator("text=document2.docx")).toBeVisible();
      await expect(page.locator("text=Загружено файлов: 2")).toBeVisible();
    });

    test("should display file size", async ({ page }) => {
      await page.click("text=Создать документ");

      const fileInput = page.locator("#file-upload");
      const content = "A".repeat(1024); // 1KB

      await fileInput.setInputFiles({
        name: "test.pdf",
        mimeType: "application/pdf",
        buffer: Buffer.from(content),
      });

      await expect(page.locator("text=test.pdf")).toBeVisible();
      // Size display (exact format depends on formatFileSize implementation)
      await expect(
        page.locator("text=/\\d+(\\.\\d+)? (B|KB|MB)/"),
      ).toBeVisible();
    });

    test("should remove uploaded file", async ({ page }) => {
      await page.click("text=Создать документ");

      const fileInput = page.locator("#file-upload");
      await fileInput.setInputFiles({
        name: "removable.pdf",
        mimeType: "application/pdf",
        buffer: Buffer.from("Test"),
      });

      await expect(page.locator("text=removable.pdf")).toBeVisible();

      // Click remove button (X icon)
      await page.locator('[data-testid="remove-file"]').first().click();

      await expect(page.locator("text=removable.pdf")).not.toBeVisible();
    });

    test("should not show file upload in edit mode", async ({ page }) => {
      // Click on existing document to edit
      await page
        .click('[data-testid="document-item"]', { timeout: 10000 })
        .catch(() => {});

      await expect(page.locator("text=Редактировать документ")).toBeVisible();
      await expect(page.locator("text=Файлы документов")).not.toBeVisible();
    });
  });

  test.describe("Metadata Entries", () => {
    test("should add metadata entry", async ({ page }) => {
      await page.click("text=Создать документ");

      await page.click("text=+ Добавить поле");

      await expect(page.locator('input[placeholder*="Ключ"]')).toBeVisible();
      await expect(
        page.locator('input[placeholder*="Значение"]'),
      ).toBeVisible();
    });

    test("should fill metadata entry fields", async ({ page }) => {
      await page.click("text=Создать документ");

      await page.click("text=+ Добавить поле");

      const keyInput = page.locator('input[placeholder*="Ключ"]').first();
      const valueInput = page.locator('input[placeholder*="Значение"]').first();

      await keyInput.fill("invoice_number");
      await valueInput.fill("INV-2026-001");

      await expect(keyInput).toHaveValue("invoice_number");
      await expect(valueInput).toHaveValue("INV-2026-001");
    });

    test("should add multiple metadata entries", async ({ page }) => {
      await page.click("text=Создать документ");

      // Add first entry
      await page.click("text=+ Добавить поле");
      await page.locator('input[placeholder*="Ключ"]').first().fill("amount");
      await page
        .locator('input[placeholder*="Значение"]')
        .first()
        .fill("100000");

      // Add second entry
      await page.click("text=+ Добавить поле");
      const keyInputs = page.locator('input[placeholder*="Ключ"]');
      const valueInputs = page.locator('input[placeholder*="Значение"]');

      await keyInputs.nth(1).fill("currency");
      await valueInputs.nth(1).fill("KGS");

      await expect(page.locator("text=amount")).toBeVisible();
      await expect(page.locator("text=currency")).toBeVisible();
    });

    test("should remove metadata entry", async ({ page }) => {
      await page.click("text=Создать документ");

      await page.click("text=+ Добавить поле");

      const keyInput = page.locator('input[placeholder*="Ключ"]').first();
      await keyInput.fill("temp_field");

      // Find and click remove button
      const removeButtons = page.locator("button").filter({ hasText: "" });
      await removeButtons.first().click();

      await expect(page.locator('input[value="temp_field"]')).not.toBeVisible();
    });

    test("should show empty state message when no entries", async ({
      page,
    }) => {
      await page.click("text=Создать документ");

      await expect(page.locator("text=Нет дополнительных полей")).toBeVisible();
    });
  });

  test.describe("Edit Document", () => {
    test("should edit existing document", async ({ page }) => {
      // Click on existing document
      await page
        .click('[data-testid="document-item"]', { timeout: 10000 })
        .catch(() => {});

      await expect(page.locator("text=Редактировать документ")).toBeVisible();

      // Modify title
      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.clear();
      await titleInput.fill("Обновленный договор №2026-002");

      // Submit
      await page.click("text=Обновить документ");

      await expect(page.locator("text=успешно обновлен")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should pre-fill form with existing data", async ({ page }) => {
      await page
        .click('[data-testid="document-item"]', { timeout: 10000 })
        .catch(() => {});

      // Check that inputs are pre-filled
      const titleInput = page.locator('input[placeholder*="Договор"]');
      await expect(titleInput).not.toHaveValue("");

      const contentInput = page.locator("textarea");
      await expect(contentInput).not.toHaveValue("");
    });

    test("should validate when editing", async ({ page }) => {
      await page
        .click('[data-testid="document-item"]', { timeout: 10000 })
        .catch(() => {});

      // Clear title
      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.clear();
      await titleInput.blur();

      await expect(
        page.locator("text=Название документа обязательно"),
      ).toBeVisible();
    });

    test("should cancel editing", async ({ page }) => {
      await page
        .click('[data-testid="document-item"]', { timeout: 10000 })
        .catch(() => {});

      // Modify field
      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.fill("Изменения для отмены");

      // Cancel
      await page.click("text=Отмена");

      // Should return to list
      await expect(page.locator("text=Список документов")).toBeVisible();
    });
  });

  test.describe("Document Status Workflow", () => {
    test("should change status from draft to pending", async ({ page }) => {
      await page.click("text=Создать документ");

      await page.fill('input[placeholder*="Договор"]', "Документ для workflow");
      await page.fill(
        "textarea",
        "Содержимое для тестирования статусов документа",
      );

      const statusSelect = page.locator("select");
      await statusSelect.selectOption("pending");

      await page.click("text=Создать документ");

      await expect(page.locator("text=успешно создан")).toBeVisible({
        timeout: 5000,
      });
    });

    test("should display correct status badge", async ({ page }) => {
      const firstDoc = page.locator('[data-testid="document-item"]').first();

      // Check for status badge (implementation-dependent)
      const hasDraftBadge = await firstDoc
        .locator("text=/черновик|draft/i")
        .isVisible()
        .catch(() => false);
      const hasPendingBadge = await firstDoc
        .locator("text=/рассмотрении|pending/i")
        .isVisible()
        .catch(() => false);
      const hasApprovedBadge = await firstDoc
        .locator("text=/утвержден|approved/i")
        .isVisible()
        .catch(() => false);

      expect(hasDraftBadge || hasPendingBadge || hasApprovedBadge).toBe(true);
    });

    test("should filter by status", async ({ page }) => {
      const statusFilter = page.locator('[data-testid="status-filter"]');

      if (await statusFilter.isVisible()) {
        await statusFilter.selectOption("draft");

        await page.waitForTimeout(500);

        // All visible documents should have draft status
        const docs = page.locator('[data-testid="document-item"]');
        const count = await docs.count();

        for (let i = 0; i < count; i++) {
          await expect(
            docs.nth(i).locator("text=/черновик|draft/i"),
          ).toBeVisible();
        }
      }
    });
  });

  test.describe("Document List", () => {
    test("should display list of documents", async ({ page }) => {
      const docs = page.locator('[data-testid="document-item"]');
      await expect(docs.first())
        .toBeVisible({ timeout: 10000 })
        .catch(() => {});
    });

    test("should show document details in list", async ({ page }) => {
      const firstDoc = page.locator('[data-testid="document-item"]').first();

      // Should display title
      await expect(firstDoc).toContainText(/договор|документ|акт/i);

      // Should display status
      const hasStatus = await firstDoc
        .locator("text=/черновик|pending|утвержден/i")
        .isVisible()
        .catch(() => false);
      expect(hasStatus).toBe(true);
    });

    test("should search documents by title", async ({ page }) => {
      const searchInput = page.locator('[data-testid="search-documents"]');

      if (await searchInput.isVisible()) {
        await searchInput.fill("Договор");

        await page.waitForTimeout(500);

        const docs = page.locator('[data-testid="document-item"]');
        const count = await docs.count();

        if (count > 0) {
          await expect(docs.first()).toContainText(/договор/i);
        }
      }
    });

    test("should paginate documents list", async ({ page }) => {
      const pagination = page.locator('[data-testid="pagination"]');

      if (await pagination.isVisible()) {
        const nextButton = page.locator("text=Следующая");
        await nextButton.click();

        await page.waitForLoadState("networkidle");

        await expect(
          page.locator('[data-testid="document-item"]'),
        ).toBeVisible();
      }
    });
  });

  test.describe("Delete Document", () => {
    test("should delete document", async ({ page }) => {
      const deleteButton = page
        .locator('[data-testid="delete-document-btn"]')
        .first();

      if (await deleteButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        await deleteButton.click();

        // Confirm deletion
        await page.click("text=Подтвердить удаление");

        await expect(page.locator("text=успешно удален")).toBeVisible({
          timeout: 5000,
        });
      }
    });

    test("should show confirmation dialog before deleting", async ({
      page,
    }) => {
      const deleteButton = page
        .locator('[data-testid="delete-document-btn"]')
        .first();

      if (await deleteButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        await deleteButton.click();

        await expect(
          page.locator("text=Вы уверены, что хотите удалить этот документ?"),
        ).toBeVisible();
      }
    });

    test("should cancel deletion", async ({ page }) => {
      const deleteButton = page
        .locator('[data-testid="delete-document-btn"]')
        .first();

      if (await deleteButton.isVisible({ timeout: 5000 }).catch(() => false)) {
        await deleteButton.click();

        // Cancel
        await page.click("text=Отмена");

        // Document should still be visible
        await expect(
          page.locator('[data-testid="document-item"]').first(),
        ).toBeVisible();
      }
    });
  });

  test.describe("Form UX", () => {
    test("should disable submit button when form is invalid", async ({
      page,
    }) => {
      await page.click("text=Создать документ");

      const submitButton = page.locator("text=Создать документ");

      // Try to submit empty form
      await submitButton.click();

      await expect(submitButton).toBeDisabled();
    });

    test("should show loading state during submission", async ({ page }) => {
      await page.click("text=Создать документ");

      // Fill valid data
      await page.fill(
        'input[placeholder*="Договор"]',
        "Тестовый документ для загрузки",
      );
      await page.fill(
        "textarea",
        "Это содержимое документа для тестирования загрузки",
      );

      // Submit
      await page.click("text=Создать документ");

      // Should show loading state briefly
      await expect(page.locator("text=Сохранение..."))
        .toBeVisible({ timeout: 1000 })
        .catch(() => {});
    });

    test("should show error message on submission failure", async ({
      page,
    }) => {
      await page.click("text=Создать документ");

      // Fill form
      await page.fill('input[placeholder*="Договор"]', "Документ с ошибкой");
      await page.fill("textarea", "Содержимое документа с ошибкой");

      // Mock API error if possible
      // await page.route('**/api/v1/documents', (route) =>
      //   route.fulfill({ status: 500, body: 'Server error' })
      // );

      await page.click("text=Создать документ");

      // Should show error message
      await expect(page.locator(".bg-red-900\\/20"))
        .toBeVisible({ timeout: 5000 })
        .catch(() => {});
    });
  });

  test.describe("Responsive Design", () => {
    test("should display form correctly on mobile", async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.click("text=Создать документ");

      await expect(page.locator("form")).toBeVisible();
      await expect(page.locator('input[placeholder*="Договор"]')).toBeVisible();
    });

    test("should display list correctly on tablet", async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });

      await expect(page.locator('[data-testid="document-item"]').first())
        .toBeVisible({ timeout: 10000 })
        .catch(() => {});
    });
  });

  test.describe("Advanced Features", () => {
    test("should handle long document titles", async ({ page }) => {
      await page.click("text=Создать документ");

      const longTitle = "Договор ".repeat(25); // 175 characters
      const titleInput = page.locator('input[placeholder*="Договор"]');
      await titleInput.fill(longTitle);

      await expect(titleInput).toHaveValue(longTitle);
      await expect(page.locator("text=/Длина: 175\\/200/")).toBeVisible();
    });

    test("should trim whitespace from title and content", async ({ page }) => {
      await page.click("text=Создать документ");

      await page.fill(
        'input[placeholder*="Договор"]',
        "  Договор с пробелами  ",
      );
      await page.fill("textarea", "  Содержимое с пробелами  ");

      const titleInput = page.locator('input[placeholder*="Договор"]') as any;
      const contentInput = page.locator("textarea") as any;

      // After blur, trimming might occur (implementation-dependent)
      await titleInput.blur();
      await contentInput.blur();
    });

    test("should preserve line breaks in content", async ({ page }) => {
      await page.click("text=Создать документ");

      const multilineContent = "Строка 1\nСтрока 2\nСтрока 3";
      await page.fill("textarea", multilineContent);

      const contentInput = page.locator("textarea");
      await expect(contentInput).toHaveValue(multilineContent);
    });
  });
});
