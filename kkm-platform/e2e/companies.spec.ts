import { test, expect } from "@playwright/test";

test.describe("Companies Management", () => {
  // Авторизация перед каждым тестом
  test.beforeEach(async ({ page }) => {
    // Переход на страницу авторизации
    await page.goto("/login");

    // Заполнение формы логина (используйте реальные тестовые данные)
    await page.fill('input[type="email"]', "test@example.com");
    await page.fill('input[type="password"]', "password123");
    await page.click('button[type="submit"]');

    // Ожидание успешной авторизации
    await page.waitForURL("/dashboard", { timeout: 5000 });
  });

  test.describe("Список компаний", () => {
    test("отображает список компаний", async ({ page }) => {
      await page.goto("/companies");

      // Проверка заголовка страницы
      await expect(page.locator("h1")).toContainText("Компании");

      // Проверка наличия кнопки создания компании
      await expect(
        page.locator('button:has-text("Создать компанию")'),
      ).toBeVisible();
    });

    test("переходит к форме создания компании", async ({ page }) => {
      await page.goto("/companies");

      // Клик по кнопке создания
      await page.click('button:has-text("Создать компанию")');

      // Проверка URL и заголовка
      await expect(page).toHaveURL(/\/companies\/create/);
      await expect(page.locator("h1")).toContainText("Создание компании");
    });
  });

  test.describe("Создание компании", () => {
    test.beforeEach(async ({ page }) => {
      await page.goto("/companies/create");
    });

    test("отображает все обязательные поля", async ({ page }) => {
      await expect(
        page.locator('label:has-text("Название компании")'),
      ).toBeVisible();
      await expect(page.locator('label:has-text("ИНН")')).toBeVisible();
      await expect(
        page.locator('label:has-text("Юридический адрес")'),
      ).toBeVisible();
      await expect(page.locator('label:has-text("Телефон")')).toBeVisible();
      await expect(page.locator('label:has-text("Email")')).toBeVisible();
    });

    test("показывает ошибки валидации при пустых полях", async ({ page }) => {
      // Попытка отправить пустую форму
      await page.click('button:has-text("Создать компанию")');

      // Проверка отображения ошибок
      await expect(
        page.locator("text=Название компании обязательно"),
      ).toBeVisible();
      await expect(page.locator("text=ИНН обязателен")).toBeVisible();
      await expect(page.locator("text=Адрес обязателен")).toBeVisible();
      await expect(page.locator("text=Телефон обязателен")).toBeVisible();
      await expect(page.locator("text=Email обязателен")).toBeVisible();
    });

    test("валидирует формат ИНН", async ({ page }) => {
      const tinInput = page.locator("input#tin");

      // Невалидный ИНН (меньше 10 цифр)
      await tinInput.fill("123456");
      await page.click('button:has-text("Создать компанию")');

      await expect(page.locator("text=/ИНН должен содержать/")).toBeVisible();

      // Валидный ИНН (10 цифр)
      await tinInput.fill("1234567890");

      // Ошибка должна исчезнуть
      await expect(
        page.locator("text=/ИНН должен содержать/"),
      ).not.toBeVisible();
    });

    test("валидирует формат email", async ({ page }) => {
      const emailInput = page.locator("input#email");

      // Невалидный email
      await emailInput.fill("invalid-email");
      await page.click('button:has-text("Создать компанию")');

      await expect(page.locator("text=Введите корректный email")).toBeVisible();

      // Валидный email
      await emailInput.fill("test@example.com");

      // Ошибка должна исчезнуть
      await expect(
        page.locator("text=Введите корректный email"),
      ).not.toBeVisible();
    });

    test("валидирует формат телефона", async ({ page }) => {
      const phoneInput = page.locator("input#phone");

      // Невалидный телефон
      await phoneInput.fill("123");
      await page.click('button:has-text("Создать компанию")');

      await expect(
        page.locator("text=/корректный номер телефона/"),
      ).toBeVisible();

      // Валидный телефон
      await phoneInput.fill("+79991234567");

      // Ошибка должна исчезнуть
      await expect(
        page.locator("text=/корректный номер телефона/"),
      ).not.toBeVisible();
    });

    test("успешно создает компанию с валидными данными", async ({ page }) => {
      // Заполнение всех обязательных полей
      await page.fill("input#name", "Test Company E2E");
      await page.fill("input#tin", "1234567890");
      await page.fill("textarea#address", "г. Москва, ул. Тестовая, д. 1");
      await page.fill("input#phone", "+79991234567");
      await page.fill("input#email", "test-e2e@example.com");

      // Заполнение опциональных полей
      await page.fill("input#kpp", "123456789");
      await page.fill("input#ogrn", "1234567890123");
      await page.fill("input#website", "https://test-company.ru");

      // Отправка формы
      await page.click('button:has-text("Создать компанию")');

      // Проверка сообщения об успехе
      await expect(page.locator("text=Компания успешно создана")).toBeVisible({
        timeout: 5000,
      });

      // Проверка редиректа на список компаний
      await expect(page).toHaveURL(/\/companies/, { timeout: 3000 });

      // Проверка, что компания появилась в списке
      await expect(page.locator("text=Test Company E2E")).toBeVisible();
    });

    test("обрабатывает ошибку при создании дубликата", async ({ page }) => {
      // Попытка создать компанию с существующим ИНН
      await page.fill("input#name", "Duplicate Company");
      await page.fill("input#tin", "1234567890"); // Используется в предыдущем тесте
      await page.fill("textarea#address", "г. Москва, ул. Тестовая, д. 1");
      await page.fill("input#phone", "+79991234567");
      await page.fill("input#email", "duplicate@example.com");

      await page.click('button:has-text("Создать компанию")');

      // Проверка отображения ошибки
      await expect(
        page.locator('[class*="red"]').locator("text=/Ошибка/"),
      ).toBeVisible({ timeout: 5000 });
    });

    test("отменяет создание и возвращается к списку", async ({ page }) => {
      await page.fill("input#name", "Test Company");

      // Клик по кнопке "Отмена"
      await page.click('button:has-text("Отмена")');

      // Проверка редиректа
      await expect(page).toHaveURL(/\/companies/);
    });
  });

  test.describe("Редактирование компании", () => {
    test("открывает форму редактирования", async ({ page }) => {
      await page.goto("/companies");

      // Клик по первой компании в списке (или конкретной тестовой компании)
      await page.click(
        'tr:has-text("Test Company E2E") button:has-text("Редактировать")',
      );

      // Проверка URL и заголовка
      await expect(page).toHaveURL(/\/companies\/.*\/edit/);
      await expect(page.locator("h1")).toContainText("Редактирование компании");
    });

    test("отображает текущие данные компании", async ({ page }) => {
      // Предполагается, что компания уже создана
      await page.goto("/companies");
      await page.click(
        'tr:has-text("Test Company E2E") button:has-text("Редактировать")',
      );

      // Проверка, что поля заполнены
      await expect(page.locator("input#name")).toHaveValue(/Test Company E2E/);
      await expect(page.locator("input#tin")).toHaveValue(/\d{10,12}/);
      await expect(page.locator("textarea#address")).not.toBeEmpty();
    });

    test("успешно обновляет данные компании", async ({ page }) => {
      await page.goto("/companies");
      await page.click(
        'tr:has-text("Test Company E2E") button:has-text("Редактировать")',
      );

      // Изменение названия
      await page.fill("input#name", "Updated Test Company E2E");

      // Изменение статуса
      await page.selectOption("select#status", "active");

      // Сохранение изменений
      await page.click('button:has-text("Сохранить изменения")');

      // Проверка сообщения об успехе
      await expect(
        page.locator("text=Компания успешно обновлена"),
      ).toBeVisible();

      // Проверка редиректа
      await expect(page).toHaveURL(/\/companies/, { timeout: 3000 });

      // Проверка обновленного названия
      await expect(page.locator("text=Updated Test Company E2E")).toBeVisible();
    });

    test("отображает поле статуса в режиме редактирования", async ({
      page,
    }) => {
      await page.goto("/companies");
      await page.click(
        'tr:has-text("Updated Test Company E2E") button:has-text("Редактировать")',
      );

      // Проверка наличия поля статуса
      await expect(page.locator('label:has-text("Статус")')).toBeVisible();
      await expect(page.locator("select#status")).toBeVisible();

      // Проверка опций статуса
      const options = await page
        .locator("select#status option")
        .allTextContents();
      expect(options).toContain("Активна");
      expect(options).toContain("Неактивна");
      expect(options).toContain("Приостановлена");
    });
  });

  test.describe("Управление участниками", () => {
    test("отображает список участников компании", async ({ page }) => {
      await page.goto("/companies");

      // Переход к деталям компании
      await page.click('tr:has-text("Updated Test Company E2E")');

      // Проверка раздела участников
      await expect(
        page.locator('h3:has-text("Участники компании")'),
      ).toBeVisible();
    });

    test("добавляет нового участника", async ({ page }) => {
      await page.goto("/companies");
      await page.click('tr:has-text("Updated Test Company E2E")');

      // Клик по кнопке добавления участника
      await page.click('button:has-text("Добавить участника")');

      // Заполнение формы
      await page.fill("input#memberEmail", "newmember@example.com");
      await page.selectOption("select#memberRole", "employee");

      // Отправка формы
      await page.click('button:has-text("Добавить")');

      // Проверка успешного добавления
      await expect(page.locator("text=/успешно добавлен/")).toBeVisible();
      await expect(page.locator("text=newmember@example.com")).toBeVisible();
    });

    test("изменяет роль участника", async ({ page }) => {
      await page.goto("/companies");
      await page.click('tr:has-text("Updated Test Company E2E")');

      // Найти участника и изменить роль
      const memberRow = page
        .locator("text=newmember@example.com")
        .locator("..");
      await memberRow.locator("select").selectOption("manager");

      // Проверка сообщения об успехе
      await expect(
        page.locator("text=/Роль участника.*обновлена/"),
      ).toBeVisible({ timeout: 3000 });
    });

    test("удаляет участника из компании", async ({ page }) => {
      await page.goto("/companies");
      await page.click('tr:has-text("Updated Test Company E2E")');

      // Найти участника и удалить
      const memberRow = page
        .locator("text=newmember@example.com")
        .locator("..");

      // Клик по кнопке удаления
      await memberRow.locator('button[title="Удалить участника"]').click();

      // Подтверждение удаления в диалоге
      page.on("dialog", (dialog) => dialog.accept());

      // Проверка, что участник удален
      await expect(page.locator("text=/удален из компании/")).toBeVisible({
        timeout: 3000,
      });
      await expect(
        page.locator("text=newmember@example.com"),
      ).not.toBeVisible();
    });
  });

  test.describe("Удаление компании", () => {
    test("успешно удаляет компанию", async ({ page }) => {
      await page.goto("/companies");

      // Клик по кнопке удаления для тестовой компании
      const companyRow = page.locator(
        'tr:has-text("Updated Test Company E2E")',
      );
      await companyRow.locator('button:has-text("Удалить")').click();

      // Подтверждение удаления
      page.on("dialog", (dialog) => dialog.accept());

      // Проверка, что компания удалена
      await expect(
        page.locator("text=Updated Test Company E2E"),
      ).not.toBeVisible({ timeout: 3000 });
    });
  });
});
