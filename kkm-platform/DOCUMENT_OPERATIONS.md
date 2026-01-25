# 📋 Реализованные операции с документами в UI

## ✅ Завершённые задачи

### 1. Расширенный компонент DocumentRow

**Файл:** [components/documents/DocumentRow.tsx](../components/documents/DocumentRow.tsx)

Добавлены новые действия для каждого документа:

- **Просмотр** (Eye icon) - просмотр всех деталей документа
- **Отправка на одобрение** (Send icon) - отправка в статусе "draft"
- **Одобрение** (CheckCircle icon) - одобрение в статусе "sent"
- **Отклонение** (XCircle icon) - отклонение в статусе "sent"
- **Меню "Ещё"** (MoreVertical icon) с дополнительными действиями:
  - Скачать документ
  - Редактировать (только в статусе "draft")
  - В архив (для всех статусов, кроме archived)
  - Удалить (только в статусе "draft")

**Динамические кнопки:**

- Кнопки автоматически отображаются/скрываются в зависимости от статуса документа
- Используются иконки из `lucide-react` для лучшей визуализации

---

### 2. Модальные окна

#### EditDocumentModal

**Файл:** [components/documents/modals/EditDocumentModal.tsx](../components/documents/modals/EditDocumentModal.tsx)

Модаль для редактирования документа:

- Редактирование названия и содержания
- Сохранение изменений через API
- Toast-уведомления об успехе/ошибке

#### ViewDocumentModal

**Файл:** [components/documents/modals/ViewDocumentModal.tsx](../components/documents/modals/ViewDocumentModal.tsx)

Модаль для просмотра документа:

- Полная информация о документе (ID, статус, автор, дата создания)
- Содержание документа
- Метаинформация (организация, последнее обновление)

#### ConfirmModal

**Файл:** [components/documents/modals/ConfirmModal.tsx](../components/documents/modals/ConfirmModal.tsx)

Универсальная модаль подтверждения:

- Подтверждение для операций: отправка, одобрение, архивирование, удаление
- Поддержка опасных операций (красный стиль для удаления)
- Обработка ошибок с отображением в модали

#### RejectDocumentModal

**Файл:** [components/documents/modals/RejectDocumentModal.tsx](../components/documents/modals/RejectDocumentModal.tsx)

Модаль для отклонения документа:

- Обязательное поле для введения причины отклонения
- Валидация (кнопка отклонения активна только при наличии текста)
- Toast-уведомление об успехе

---

### 3. Обновлённый DocumentsList

**Файл:** [components/documents/DocumentsList.tsx](../components/documents/DocumentsList.tsx)

Главный компонент списка документов теперь включает:

**Управление состоянием:**

- `viewingDocument` - для модали просмотра
- `editingDocument` - для модали редактирования
- `rejectingDocument` - для модали отклонения
- `operationModal` - для модалей подтверждения (send, approve, archive, delete)

**Обработчики операций:**

- `handleSendDocument()` - отправка на одобрение
- `handleApproveDocument()` - одобрение
- `handleArchiveDocument()` - архивирование
- `handleDeleteDocument()` - удаление
- `handleRefresh()` - обновление списка документов

**Передача обработчиков в DocumentRow:**

```typescript
<DocumentRow
  document={doc}
  onView={(doc) => setViewingDocument(doc)}
  onEdit={(doc) => setEditingDocument(doc)}
  onDelete={(doc) => setOperationModal({ type: 'delete', document: doc })}
  onArchive={(doc) => setOperationModal({ type: 'archive', document: doc })}
  onSend={(doc) => setOperationModal({ type: 'send', document: doc })}
  onApprove={(doc) => setOperationModal({ type: 'approve', document: doc })}
  onReject={(doc) => setRejectingDocument(doc)}
  onDownload={(doc) => console.log('Загрузка документа:', doc)}
/>
```

---

### 4. Toast-уведомления

**Интеграция с:** `react-hot-toast`

Используется существующая система уведомлений:

**Успешные операции:**

```typescript
toast.success(`Документ "${doc.title}" отправлен на одобрение`);
toast.success(`Документ "${doc.title}" одобрен`);
toast.success(`Документ "${doc.title}" архивирован`);
toast.success(`Документ "${doc.title}" удалён`);
toast.success("Документ успешно обновлен");
toast.success("Документ отклонен");
```

**Ошибки:**

```typescript
toast.error(errorMessage);
```

---

## 🔌 Используемые API функции

Все функции из [lib/api/documents.ts](../lib/api/documents.ts):

| Операция      | Функция             | Метод | Эндпоинт                  |
| ------------- | ------------------- | ----- | ------------------------- |
| Обновление    | `updateDocument()`  | PUT   | `/documents/{id}`         |
| Отправка      | `sendDocument()`    | POST  | `/documents/{id}/send`    |
| Одобрение     | `approveDocument()` | POST  | `/documents/{id}/approve` |
| Отклонение    | `rejectDocument()`  | POST  | `/documents/{id}/reject`  |
| Архивирование | `archiveDocument()` | POST  | `/documents/{id}/archive` |

---

## 📁 Структура файлов

```
components/
  ├── documents/
  │   ├── DocumentRow.tsx          (✨ расширенный)
  │   ├── DocumentsList.tsx         (✨ обновлённый)
  │   ├── CreateDocumentForm.tsx
  │   └── modals/
  │       ├── index.ts             (✨ новый)
  │       ├── EditDocumentModal.tsx    (✨ новый)
  │       ├── ViewDocumentModal.tsx    (✨ новый)
  │       ├── ConfirmModal.tsx         (✨ новый)
  │       └── RejectDocumentModal.tsx  (✨ новый)
```

---

## 🎨 Визуальные особенности

### Иконки операций

- 👁️ **Eye** - Просмотр
- 📤 **Send** - Отправить
- ✅ **CheckCircle** - Одобрить (зелёная)
- ❌ **XCircle** - Отклонить (красная)
- ✏️ **Edit** - Редактировать
- 🗑️ **Trash2** - Удалить (красная)
- 📦 **Archive** - В архив
- ⋮ **MoreVertical** - Меню

### Цветовая схема

- Зелёная (Success): Одобрение, отправка
- Красная (Danger): Отклонение, удаление
- Синяя (Info): Редактирование
- Серая (Default): Архивирование

---

## 🚀 Тестирование

Все операции готовы к тестированию:

1. **Просмотр документа** - нажать на Eye
2. **Редактирование** - нажать на ⋮ → Редактировать
3. **Отправка** - нажать на Send (видна только для draft)
4. **Одобрение** - нажать на ✅ (видна только для sent)
5. **Отклонение** - нажать на ❌ (видна только для sent)
6. **Архивирование** - нажать на ⋮ → В архив
7. **Удаление** - нажать на ⋮ → Удалить (видна только для draft)

---

## ✅ Статус

- ✅ Frontend компилируется успешно
- ✅ Все API функции интегрированы
- ✅ Модальные окна готовы к использованию
- ✅ Toast-уведомления настроены
- ✅ Динамическая видимость кнопок по статусу документа
