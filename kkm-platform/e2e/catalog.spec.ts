import { test, expect } from "@playwright/test";

test.describe("Catalog Management", () => {
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

  test.describe("Список каталога", () => {
    test("отображает список товаров", async ({ page }) => {
      await page.goto("/catalog");

      // Проверка заголовка страницы
      await expect(page.locator("h1")).toContainText("Каталог");

      // Проверка наличия кнопки создания товара
      await expect(
        page.locator('button:has-text("Добавить товар")'),
      ).toBeVisible();
    });

    test("переходит к форме создания товара", async ({ page }) => {
      await page.goto("/catalog");

      // Клик по кнопке создания
      await page.click('button:has-text("Добавить товар")');

      // Проверка URL и заголовка
      await expect(page).toHaveURL(/\/catalog\/create/);
      await expect(page.locator("h1")).toContainText(
        "Добавление товара в каталог",
      );
    });
  });

  test.describe("Создание товара", () => {
    test.beforeEach(async ({ page }) => {
      await page.goto("/catalog/create");
    });

    test("отображает все обязательные поля", async ({ page }) => {
      await expect(
        page.locator('label:has-text("Наименование товара")'),
      ).toBeVisible();
      await expect(
        page.locator('label:has-text("Номер по каталогу")'),
      ).toBeVisible();
      await expect(page.locator('label:has-text("Код ТНВЭД")')).toBeVisible();
      await expect(page.locator('label:has-text("Цена")')).toBeVisible();
      await expect(page.locator('label:has-text("Валюта")')).toBeVisible();
      await expect(
        page.locator('label:has-text("Единица измерения")'),
      ).toBeVisible();
    });

    test("показывает ошибки валидации при пустых полях", async ({ page }) => {
      // Попытка отправить пустую форму
      await page.click('button:has-text("Добавить в каталог")');

      // Проверка отображения ошибок
      await expect(
        page.locator("text=Название товара/услуги обязательно"),
      ).toBeVisible();
      await expect(
        page.locator("text=Номер по каталогу обязателен"),
      ).toBeVisible();
      await expect(
        page.locator("text=Код ТНВЭД обязателен для ЭСФ"),
      ).toBeVisible();
      await expect(page.locator("text=Цена обязательна")).toBeVisible();
    });

    test("валидирует формат ТНВЭД кода", async ({ page }) => {
      const tnvedInput = page.locator("input#tnved_code");

      // Невалидный ТНВЭД (меньше 10 цифр)
      await tnvedInput.fill("12345");
      await page.click('button:has-text("Добавить в каталог")');

      await expect(
        page.locator("text=Код ТНВЭД должен содержать 10 цифр"),
      ).toBeVisible();

      // Валидный ТНВЭД (10 цифр)
      await tnvedInput.fill("1234567890");

      // Ошибка должна исчезнуть и код должен быть отформатирован
      await expect(
        page.locator("text=Код ТНВЭД должен содержать 10 цифр"),
      ).not.toBeVisible();
      await expect(tnvedInput).toHaveValue("1234 56 7890");
    });

    test("автоматически форматирует ТНВЭД код при вводе", async ({ page }) => {
      const tnvedInput = page.locator("input#tnved_code");

      // Ввод без пробелов
      await tnvedInput.fill("1234567890");

      // Проверка автоформатирования (XXXX XX XXXX)
      await expect(tnvedInput).toHaveValue("1234 56 7890");
    });

    test("валидирует цену (положительное число)", async ({ page }) => {
      const priceInput = page.locator("input#price");

      // Пустая цена
      await page.click('button:has-text("Добавить в каталог")');
      await expect(page.locator("text=Цена обязательна")).toBeVisible();

      // Валидная цена
      await priceInput.fill("99.99");
      await expect(
        page.locator("text=Цена должна быть положительным числом"),
      ).not.toBeVisible();
    });

    test("успешно создает товар с валидными данными", async ({ page }) => {
      // Заполнение всех обязательных полей
      await page.fill("input#name", "Test Product E2E");
      await page.fill("input#number", "PROD-E2E-001");
      await page.fill("input#tnved_code", "1234567890");
      await page.fill("input#price", "150.50");

      // Валюта и единица измерения уже установлены по умолчанию (KGS, 796)

      // Опциональное описание
      await page.fill(
        "textarea#description",
        "This is a test product for E2E testing",
      );

      // Отправка формы
      await page.click('button:has-text("Добавить в каталог")');

      // Проверка сообщения об успехе
      await expect(
        page.locator("text=Товар успешно добавлен в каталог"),
      ).toBeVisible({ timeout: 5000 });

      // Проверка редиректа на список каталога
      await expect(page).toHaveURL(/\/catalog/, { timeout: 3000 });

      // Проверка, что товар появился в списке
      await expect(page.locator("text=Test Product E2E")).toBeVisible();
    });

    test("выбирает валюту из справочника", async ({ page }) => {
      const currencySelect = page.locator("select#currency");

      // Проверка наличия основных валют
      await expect(currencySelect.locator('option[value="KGS"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="USD"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="RUB"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="EUR"]')).toBeVisible();

      // Выбор USD
      await currencySelect.selectOption("USD");

      // Проверка отображения символа $
      await expect(page.locator("text=$").first()).toBeVisible();
    });

    test("выбирает единицу измерения из справочника ОКЕИ", async ({ page }) => {
      const unitSelect = page.locator("select#unit");

      // Проверка наличия основных единиц
      await expect(unitSelect.locator('option[value="796"]')).toBeVisible(); // шт
      await expect(unitSelect.locator('option[value="166"]')).toBeVisible(); // кг
      await expect(unitSelect.locator('option[value="112"]')).toBeVisible(); // л

      // Выбор "кг" (166)
      await unitSelect.selectOption("166");

      // Проверка отображения выбранной единицы
      await expect(page.locator("text=Выбранная единица: кг")).toBeVisible();
    });

    test("обрабатывает ошибку при создании дубликата", async ({ page }) => {
      // Попытка создать товар с существующим номером
      await page.fill("input#name", "Duplicate Product");
      await page.fill("input#number", "PROD-E2E-001"); // Используется в предыдущем тесте
      await page.fill("input#tnved_code", "9876543210");
      await page.fill("input#price", "100");

      await page.click('button:has-text("Добавить в каталог")');

      // Проверка отображения ошибки
      await expect(
        page.locator('[class*="red"]').locator("text=/Ошибка/"),
      ).toBeVisible({ timeout: 5000 });
    });

    test("отменяет создание и возвращается к списку", async ({ page }) => {
      await page.fill("input#name", "Test Product");

      // Клик по кнопке "Отмена"
      await page.click('button:has-text("Отмена")');

      // Проверка редиректа
      await expect(page).toHaveURL(/\/catalog/);
    });
  });

  test.describe("Редактирование товара", () => {
    test("открывает форму редактирования", async ({ page }) => {
      await page.goto("/catalog");

      // Клик по кнопке редактирования для тестового товара
      await page.click(
        'tr:has-text("Test Product E2E") button:has-text("Редактировать")',
      );

      // Проверка URL и заголовка
      await expect(page).toHaveURL(/\/catalog\/.*\/edit/);
      await expect(page.locator("h1")).toContainText("Редактирование товара");
    });

    test("отображает текущие данные товара", async ({ page }) => {
      await page.goto("/catalog");
      await page.click(
        'tr:has-text("Test Product E2E") button:has-text("Редактировать")',
      );

      // Проверка, что поля заполнены
      await expect(page.locator("input#name")).toHaveValue(/Test Product E2E/);
      await expect(page.locator("input#number")).toHaveValue(/PROD-E2E-001/);
      await expect(page.locator("input#tnved_code")).toHaveValue(
        /1234 56 7890/,
      ); // Форматированный
      await expect(page.locator("input#price")).toHaveValue(/150/);
    });

    test("успешно обновляет данные товара", async ({ page }) => {
      await page.goto("/catalog");
      await page.click(
        'tr:has-text("Test Product E2E") button:has-text("Редактировать")',
      );

      // Изменение названия и цены
      await page.fill("input#name", "Updated Test Product E2E");
      await page.fill("input#price", "200.00");

      // Сохранение изменений
      await page.click('button:has-text("Сохранить изменения")');

      // Проверка сообщения об успехе
      await expect(page.locator("text=Товар успешно обновлен")).toBeVisible();

      // Проверка редиректа
      await expect(page).toHaveURL(/\/catalog/, { timeout: 3000 });

      // Проверка обновленного названия
      await expect(page.locator("text=Updated Test Product E2E")).toBeVisible();
    });

    test("изменяет валюту и единицу измерения", async ({ page }) => {
      await page.goto("/catalog");
      await page.click(
        'tr:has-text("Updated Test Product E2E") button:has-text("Редактировать")',
      );

      // Изменение валюты на USD
      await page.selectOption("select#currency", "USD");

      // Изменение единицы измерения на кг
      await page.selectOption("select#unit", "166");

      // Сохранение изменений
      await page.click('button:has-text("Сохранить изменения")');

      await expect(page.locator("text=Товар успешно обновлен")).toBeVisible();
    });

    test("обновляет ТНВЭД код с автоформатированием", async ({ page }) => {
      await page.goto("/catalog");
      await page.click(
        'tr:has-text("Updated Test Product E2E") button:has-text("Редактировать")',
      );

      const tnvedInput = page.locator("input#tnved_code");

      // Очистка и ввод нового ТНВЭД кода
      await tnvedInput.clear();
      await tnvedInput.fill("9876543210");

      // Проверка автоформатирования
      await expect(tnvedInput).toHaveValue("9876 54 3210");

      // Сохранение
      await page.click('button:has-text("Сохранить изменения")');

      await expect(page.locator("text=Товар успешно обновлен")).toBeVisible();
    });
  });

  test.describe("Поиск и фильтрация", () => {
    test("ищет товар по названию", async ({ page }) => {
      await page.goto("/catalog");

      // Ввод в поле поиска
      const searchInput = page.locator('input[placeholder*="Поиск"]');
      await searchInput.fill("Updated Test Product");

      // Проверка, что найден нужный товар
      await expect(page.locator("text=Updated Test Product E2E")).toBeVisible();
    });

    test("фильтрует товары по валюте", async ({ page }) => {
      await page.goto("/catalog");

      // Если есть фильтр по валюте
      const currencyFilter = page.locator(
        'select[aria-label="Фильтр по валюте"]',
      );
      if (await currencyFilter.isVisible()) {
        await currencyFilter.selectOption("USD");

        // Проверка, что отображаются только товары в USD
        await expect(page.locator('td:has-text("$")')).toBeVisible();
      }
    });
  });

  test.describe("Удаление товара", () => {
    test("успешно удаляет товар", async ({ page }) => {
      await page.goto("/catalog");

      // Клик по кнопке удаления для тестового товара
      const productRow = page.locator(
        'tr:has-text("Updated Test Product E2E")',
      );
      await productRow.locator('button:has-text("Удалить")').click();

      // Подтверждение удаления
      page.on("dialog", (dialog) => dialog.accept());

      // Проверка, что товар удален
      await expect(
        page.locator("text=Updated Test Product E2E"),
      ).not.toBeVisible({ timeout: 3000 });
    });
  });

  test.describe("Проверка справочников", () => {
    test("отображает полный список валют", async ({ page }) => {
      await page.goto("/catalog/create");

      const currencySelect = page.locator("select#currency");

      // Проверка всех валют из CURRENCIES
      await expect(currencySelect.locator('option[value="KGS"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="USD"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="RUB"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="EUR"]')).toBeVisible();
      await expect(currencySelect.locator('option[value="CNY"]')).toBeVisible();
    });

    test("отображает полный список единиц измерения ОКЕИ", async ({ page }) => {
      await page.goto("/catalog/create");

      const unitSelect = page.locator("select#unit");

      // Проверка основных единиц из UNITS_OF_MEASURE
      await expect(unitSelect.locator('option[value="796"]')).toBeVisible(); // шт
      await expect(unitSelect.locator('option[value="006"]')).toBeVisible(); // м
      await expect(unitSelect.locator('option[value="166"]')).toBeVisible(); // кг
      await expect(unitSelect.locator('option[value="112"]')).toBeVisible(); // л
      await expect(unitSelect.locator('option[value="212"]')).toBeVisible(); // т
    });
  });

  test.describe("Обработка ошибок", () => {
    test("показывает ошибку при недоступном API", async ({ page }) => {
      // Имитация недоступности API (требует настройки mock)
      // Этот тест может быть расширен с использованием route mocking
    });

    test("показывает ошибку при превышении лимита цены", async ({ page }) => {
      await page.goto("/catalog/create");

      await page.fill("input#name", "Expensive Product");
      await page.fill("input#number", "PROD-EXP-001");
      await page.fill("input#tnved_code", "1111111111");
      await page.fill("input#price", "999999999999"); // Очень большая цена

      await page.click('button:has-text("Добавить в каталог")');

      // Может быть ошибка валидации на backend
      // Проверка обработки
    });
  });
});
