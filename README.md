# Expenses registration

This program registers expenses

## Data model

```mermaid
erDiagram
    %% Relationships
    CATEGORY ||--o{ EXPENSE : "categorizes"
    PAYMENT_METHOD ||--o{ EXPENSE : "funds"
    CATEGORY ||--o{ CATEGORY : "parent of"
    %% Many-to-many for tags
    EXPENSE }o--o{ TAG : "has"
    %% (or explicitly with junction table - see note below)

    EXPENSE {
        int expense_id PK
        int category_id FK
        int payment_method_id FK
        char(3) currency "default 'PEN'"
        %% Store monetary amounts as integers (not decimal) to avoid floating point rounding errors
        bigint amount
        timestamptz expense_date
        string merchant_name
        string description
        datetime created_at
        datetime updated_at
        datetime deleted_at "soft delete"
        %% CHECK (amount > 0)
    }

    TAG {
        int tag_id PK
        string tag_name "unique (business, japan_trip_2026, groceries...)"
        string color "hex #RRGGBB (optional)"
        string icon "icon name"
        datetime created_at
        datetime updated_at
    }

    %% Junction table (explicit many-to-many)
    EXPENSE_TAG {
        int expense_id FK
        int tag_id FK
        datetime created_at "when this tag was added"
    }

    CATEGORY {
        int category_id PK
        int parent_category_id FK "nullable → hierarchy"
        string category_name
        string icon "icon name"
        string color "hex #RRGGBB"
        datetime created_at
        datetime updated_at
        %% UNIQUE CONSTRAINT (parent_category_id, category_name) to allow having “Food” under both “Personal” and “Business” without collision
    }

    PAYMENT_METHOD {
        int payment_method_id PK
        string method_name
        enum method_type "credit_card,cash,debit_card,yape,plin,bank_transfer"
        string icon
        datetime created_at
        datetime updated_at
    }
```

## Icons used in this project

This project will use an open source icon library in the future, so the fields named icon should be string.

## Technologies used in this project

- PostgreSQL 18
- Go 1.26.2
- Echo v5
