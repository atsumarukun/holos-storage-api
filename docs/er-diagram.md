```mermaid
erDiagram

volumes {
  char(36) id PK
  char(36) account_id
  varchar(255) name
  tinyint(1) is_public
  datetime(6) created_at
  datetime(6) updated_at
  datetime(6) deleted_at
}
```
