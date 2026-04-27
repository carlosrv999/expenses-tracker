CREATE TYPE payment_method_type AS ENUM (
    'credit_card',
    'cash',
    'debit_card',
    'yape',
    'plin',
    'bank_transfer'
);

CREATE TABLE category (
    category_id        SERIAL PRIMARY KEY,
    parent_category_id INTEGER REFERENCES category(category_id) ON DELETE SET NULL,
    category_name      TEXT    NOT NULL,
    icon               TEXT,
    color              CHAR(7),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE NULLS NOT DISTINCT (parent_category_id, category_name)
);

CREATE TABLE payment_method (
    payment_method_id SERIAL PRIMARY KEY,
    method_name       TEXT                NOT NULL UNIQUE,
    method_type       payment_method_type NOT NULL,
    icon              TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tag (
    tag_id     SERIAL PRIMARY KEY,
    tag_name   TEXT NOT NULL UNIQUE,
    color      CHAR(7),
    icon       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE expense (
    expense_id        SERIAL PRIMARY KEY,
    category_id       INTEGER     NOT NULL REFERENCES category(category_id)             ON DELETE RESTRICT,
    payment_method_id INTEGER     NOT NULL REFERENCES payment_method(payment_method_id) ON DELETE RESTRICT,
    currency          CHAR(3)     NOT NULL DEFAULT 'PEN',
    amount            BIGINT      NOT NULL CHECK (amount > 0),
    expense_date      TIMESTAMPTZ NOT NULL,
    merchant_name     TEXT,
    description       TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_expense_category_id       ON expense(category_id);
CREATE INDEX idx_expense_payment_method_id ON expense(payment_method_id);
CREATE INDEX idx_expense_expense_date      ON expense(expense_date);
CREATE INDEX idx_expense_deleted_at        ON expense(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE expense_tag (
    expense_id INTEGER     NOT NULL REFERENCES expense(expense_id) ON DELETE CASCADE,
    tag_id     INTEGER     NOT NULL REFERENCES tag(tag_id)         ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (expense_id, tag_id)
);

CREATE INDEX idx_expense_tag_tag_id ON expense_tag(tag_id);
