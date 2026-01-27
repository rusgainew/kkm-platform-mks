# Page snapshot

```yaml
- generic [active] [ref=e1]:
  - generic [ref=e4]:
    - generic [ref=e5]:
      - heading "Счета" [level=2] [ref=e6]:
        - img [ref=e7]
        - text: Счета
      - link "Создать счет" [ref=e10] [cursor=pointer]:
        - /url: /invoices/create
        - img [ref=e11]
        - text: Создать счет
    - textbox "Поиск по номеру или компании..." [ref=e13]
    - generic [ref=e14]:
      - table [ref=e15]:
        - rowgroup [ref=e16]:
          - row "Номер Компания Дата Сумма Статус Действия" [ref=e17]:
            - columnheader "Номер" [ref=e18]
            - columnheader "Компания" [ref=e19]
            - columnheader "Дата" [ref=e20]
            - columnheader "Сумма" [ref=e21]
            - columnheader "Статус" [ref=e22]
            - columnheader "Действия" [ref=e23]
        - rowgroup
      - generic [ref=e24]: Нет счетов
  - generic [ref=e25]:
    - img [ref=e27]
    - button "Open Tanstack query devtools" [ref=e75] [cursor=pointer]:
      - img [ref=e76]
  - button "Open Next.js Dev Tools" [ref=e129] [cursor=pointer]:
    - img [ref=e130]
  - alert [ref=e133]
```