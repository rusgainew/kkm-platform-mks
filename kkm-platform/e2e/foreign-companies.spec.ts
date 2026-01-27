import { test, expect, type Page } from "@playwright/test";

const TEST_EMAIL = "admin@example.com";
const TEST_PASSWORD = "admin123";
const BASE_URL =
  process.env.PLAYWRIGHT_TEST_BASE_URL || "http://localhost:3000";

/**
 * Хелпер для авторизации
 */
async function login(page: Page) {
  await page.goto(`${BASE_URL}/login`);
  await page.fill('input[name="email"]', TEST_EMAIL);
  await page.fill('input[name="password"]', TEST_PASSWORD);
  await page.click('button[type="submit"]');
  await page.waitForURL("**/dashboard", { timeout: 10000 });
}

/**
 * Хелпер для генерации случайного Tax ID
 */
function generateRandomTaxId(): string {
  return `TAX-${Math.random().toString(36).substring(2, 11).toUpperCase()}`;
}

test.describe("Foreign Companies E2E Tests", () => {
  test.beforeEach(async ({ page }) => {
    await login(page);
    // Переходим на страницу иностранных компаний
    await page.goto(`${BASE_URL}/foreign-companies`);
  });

  test.describe("Create Foreign Company", () => {
    test("should create a new foreign company with valid data", async ({
      page,
    }) => {
      // Нажимаем кнопку "Добавить компанию"
      await page.click('button:has-text("Добавить компанию")');

      // Ждем появления формы
      await expect(
        page.locator('text="Добавить иностранную компанию"'),
      ).toBeVisible();

      // Заполняем форму
      const taxId = generateRandomTaxId();
      await page.fill('input[name="name"]', "Test Foreign Company LLC");
      await page.fill('input[name="tax_id"]', taxId);
      await page.selectOption('select[name="country"]', "US");
      await page.fill(
        'textarea[name="address"]',
        "123 Main Street, Suite 100, New York, NY 10001, USA",
      );
      await page.fill('input[name="contact_email"]', "contact@testforeign.com");
      await page.fill('input[name="contact_phone"]', "+1 (555) 987-6543");
      await page.selectOption('select[name="currency"]', "USD");

      // Отправляем форму
      await page.click('button[type="submit"]:has-text("Создать компанию")');

      // Проверяем успешное создание
      await expect(page.locator('text="успешно создана"')).toBeVisible({
        timeout: 5000,
      });

      // Проверяем, что компания появилась в списке
      await expect(
        page.locator('text="Test Foreign Company LLC"'),
      ).toBeVisible();
    });

    test("should show validation errors for empty required fields", async ({
      page,
    }) => {
      await page.click('button:has-text("Добавить компанию")');

      // Нажимаем кнопку отправки без заполнения полей
      await page.click('button[type="submit"]:has-text("Создать компанию")');

      // Проверяем наличие ошибок валидации
      await expect(
        page.locator('text="Название компании обязательно"'),
      ).toBeVisible();
      await expect(
        page.locator('text="Налоговый номер (PIN/TIN) обязателен"'),
      ).toBeVisible();
      await expect(page.locator('text="Код страны обязателен"')).toBeVisible();
      await expect(page.locator('text="Адрес обязателен"')).toBeVisible();
      await expect(page.locator('text="Email обязателен"')).toBeVisible();
      await expect(page.locator('text="Телефон обязателен"')).toBeVisible();
    });

    test("should show error for invalid email format", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      await page.fill('input[name="name"]', "Test Company");
      await page.fill('input[name="tax_id"]', "TAX-12345");
      await page.selectOption('select[name="country"]', "US");
      await page.fill('textarea[name="address"]', "123 Main Street");
      await page.fill('input[name="contact_email"]', "invalid-email"); // Invalid
      await page.fill('input[name="contact_phone"]', "+1 (555) 123-4567");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(page.locator('text="Неверный формат email"')).toBeVisible();
    });

    test("should show error for invalid phone format", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      await page.fill('input[name="name"]', "Test Company");
      await page.fill('input[name="tax_id"]', "TAX-12345");
      await page.selectOption('select[name="country"]', "US");
      await page.fill('textarea[name="address"]', "123 Main Street");
      await page.fill('input[name="contact_email"]', "test@example.com");
      await page.fill('input[name="contact_phone"]', "123"); // Too short

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(
        page.locator('text="Неверный формат телефона"'),
      ).toBeVisible();
    });

    test("should normalize tax_id to uppercase", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      const taxId = "abc-123-def";
      await page.fill('input[name="name"]', "Test Company");
      await page.fill('input[name="tax_id"]', taxId);
      await page.selectOption('select[name="country"]', "CN");
      await page.fill('textarea[name="address"]', "123 Beijing Road, Beijing");
      await page.fill('input[name="contact_email"]', "test@example.com");
      await page.fill('input[name="contact_phone"]', "+86 10 1234 5678");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(page.locator('text="успешно создана"')).toBeVisible({
        timeout: 5000,
      });

      // Tax ID должен быть в верхнем регистре
      await expect(page.locator('text="ABC-123-DEF"')).toBeVisible();
    });

    test("should create company with popular country (Russia)", async ({
      page,
    }) => {
      await page.click('button:has-text("Добавить компанию")');

      const taxId = generateRandomTaxId();
      await page.fill('input[name="name"]', "Russian Import Company");
      await page.fill('input[name="tax_id"]', taxId);
      await page.selectOption('select[name="country"]', "RU");
      await page.fill(
        'textarea[name="address"]',
        "ул. Тверская, д. 1, Москва, 125009",
      );
      await page.fill('input[name="contact_email"]', "info@rusimport.ru");
      await page.fill('input[name="contact_phone"]', "+7 495 123 4567");
      await page.selectOption('select[name="currency"]', "RUB");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(page.locator('text="успешно создана"')).toBeVisible({
        timeout: 5000,
      });
      await expect(page.locator('text="Russian Import Company"')).toBeVisible();
    });
  });

  test.describe("Edit Foreign Company", () => {
    test("should edit existing foreign company", async ({ page }) => {
      // Предполагаем, что в списке уже есть компании
      // Кликаем на первую компанию для редактирования
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      // Ждем загрузки формы редактирования
      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      // Изменяем название
      const newName = `Updated Company ${Date.now()}`;
      await page.fill('input[name="name"]', newName);

      // Сохраняем изменения
      await page.click('button[type="submit"]:has-text("Сохранить изменения")');

      // Проверяем успешное обновление
      await expect(page.locator('text="успешно обновлена"')).toBeVisible({
        timeout: 5000,
      });

      // Проверяем, что новое название отображается
      await expect(page.locator(`text="${newName}"`)).toBeVisible();
    });

    test("should update company contact information", async ({ page }) => {
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      // Меняем email и телефон
      await page.fill('input[name="contact_email"]', "newemail@example.com");
      await page.fill('input[name="contact_phone"]', "+1 (555) 999-8888");

      await page.click('button[type="submit"]:has-text("Сохранить изменения")');

      await expect(page.locator('text="успешно обновлена"')).toBeVisible({
        timeout: 5000,
      });
    });

    test("should update company address", async ({ page }) => {
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      const newAddress = "456 New Street, Los Angeles, CA 90001, USA";
      await page.fill('textarea[name="address"]', newAddress);

      await page.click('button[type="submit"]:has-text("Сохранить изменения")');

      await expect(page.locator('text="успешно обновлена"')).toBeVisible({
        timeout: 5000,
      });
    });

    test("should update company currency", async ({ page }) => {
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      await page.selectOption('select[name="currency"]', "EUR");

      await page.click('button[type="submit"]:has-text("Сохранить изменения")');

      await expect(page.locator('text="успешно обновлена"')).toBeVisible({
        timeout: 5000,
      });
    });
  });

  test.describe("Delete Foreign Company", () => {
    test("should delete foreign company with confirmation", async ({
      page,
    }) => {
      // Получаем название первой компании для проверки
      const firstCompanyName = await page
        .locator('[data-testid="company-name"]:first-of-type')
        .textContent();

      // Кликаем на кнопку редактирования
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      // Нажимаем кнопку удаления
      await page.click('button:has-text("Удалить")');

      // Проверяем появление модального окна подтверждения
      await expect(page.locator('text="Подтвердите удаление"')).toBeVisible();

      // Подтверждаем удаление
      await page.click('button:has-text("Удалить"):last-of-type');

      // Проверяем успешное удаление
      await expect(page.locator('text="успешно удалена"')).toBeVisible({
        timeout: 5000,
      });

      // Проверяем, что компания больше не в списке
      if (firstCompanyName) {
        await expect(
          page.locator(`text="${firstCompanyName}"`),
        ).not.toBeVisible();
      }
    });

    test("should cancel deletion on cancel button", async ({ page }) => {
      await page.click('button[aria-label="Редактировать"]:first-of-type');

      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();

      // Нажимаем кнопку удаления
      await page.click('button:has-text("Удалить")');

      // Проверяем появление модального окна
      await expect(page.locator('text="Подтвердите удаление"')).toBeVisible();

      // Нажимаем отмену
      await page.click('button:has-text("Отмена"):last-of-type');

      // Модальное окно должно закрыться
      await expect(
        page.locator('text="Подтвердите удаление"'),
      ).not.toBeVisible();

      // Форма редактирования все еще видна
      await expect(
        page.locator('text="Редактировать иностранную компанию"'),
      ).toBeVisible();
    });
  });

  test.describe("Foreign Companies List", () => {
    test("should display list of foreign companies", async ({ page }) => {
      // Проверяем наличие списка
      await expect(page.locator('[data-testid="companies-list"]')).toBeVisible({
        timeout: 10000,
      });
    });

    test("should search foreign companies by name", async ({ page }) => {
      // Вводим поисковый запрос
      await page.fill('input[placeholder*="Поиск"]', "Test");

      // Ждем результатов поиска
      await page.waitForTimeout(500);

      // Проверяем, что результаты содержат "Test"
      const results = page.locator('[data-testid="company-name"]');
      const count = await results.count();

      if (count > 0) {
        for (let i = 0; i < count; i++) {
          const text = await results.nth(i).textContent();
          expect(text?.toLowerCase()).toContain("test");
        }
      }
    });

    test("should filter companies by country", async ({ page }) => {
      // Выбираем фильтр по стране
      await page.selectOption('select[name="country_filter"]', "US");

      // Ждем применения фильтра
      await page.waitForTimeout(500);

      // Проверяем, что все результаты имеют код страны US
      const countryBadges = page.locator('[data-testid="company-country"]');
      const count = await countryBadges.count();

      if (count > 0) {
        for (let i = 0; i < count; i++) {
          const text = await countryBadges.nth(i).textContent();
          expect(text).toContain("US");
        }
      }
    });

    test("should paginate through companies list", async ({ page }) => {
      // Проверяем наличие пагинации
      const nextButton = page.locator('button:has-text("Следующая")');

      if (await nextButton.isVisible()) {
        // Кликаем на следующую страницу
        await nextButton.click();

        // Ждем загрузки новых данных
        await page.waitForTimeout(500);

        // Проверяем, что URL изменился (содержит page=2)
        expect(page.url()).toContain("page=2");
      }
    });

    test("should display company details on click", async ({ page }) => {
      // Кликаем на первую компанию
      await page.click('[data-testid="company-row"]:first-of-type');

      // Проверяем отображение деталей компании
      await expect(page.locator('[data-testid="company-details"]')).toBeVisible(
        {
          timeout: 5000,
        },
      );

      // Проверяем наличие ключевых полей
      await expect(page.locator('text="Налоговый номер"')).toBeVisible();
      await expect(page.locator('text="Страна"')).toBeVisible();
      await expect(page.locator('text="Адрес"')).toBeVisible();
      await expect(page.locator('text="Email"')).toBeVisible();
      await expect(page.locator('text="Телефон"')).toBeVisible();
    });
  });

  test.describe("Form Validation", () => {
    test("should validate tax ID length", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      await page.fill('input[name="name"]', "Test Company");
      await page.fill('input[name="tax_id"]', "123"); // Too short
      await page.selectOption('select[name="country"]', "US");
      await page.fill('textarea[name="address"]', "123 Main Street");
      await page.fill('input[name="contact_email"]', "test@example.com");
      await page.fill('input[name="contact_phone"]', "+1 (555) 123-4567");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(
        page.locator('text="Неверный формат налогового номера"'),
      ).toBeVisible();
    });

    test("should validate address length", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      await page.fill('input[name="name"]', "Test Company");
      await page.fill('input[name="tax_id"]', "TAX-12345");
      await page.selectOption('select[name="country"]', "US");
      await page.fill('textarea[name="address"]', "123"); // Too short
      await page.fill('input[name="contact_email"]', "test@example.com");
      await page.fill('input[name="contact_phone"]', "+1 (555) 123-4567");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(
        page.locator('text="Адрес должен содержать минимум 5 символов"'),
      ).toBeVisible();
    });

    test("should validate company name length", async ({ page }) => {
      await page.click('button:has-text("Добавить компанию")');

      await page.fill('input[name="name"]', "A"); // Too short
      await page.fill('input[name="tax_id"]', "TAX-12345");
      await page.selectOption('select[name="country"]', "US");
      await page.fill('textarea[name="address"]', "123 Main Street");
      await page.fill('input[name="contact_email"]', "test@example.com");
      await page.fill('input[name="contact_phone"]', "+1 (555) 123-4567");

      await page.click('button[type="submit"]:has-text("Создать компанию")');

      await expect(
        page.locator('text="Название должно содержать минимум 2 символа"'),
      ).toBeVisible();
    });
  });
});
