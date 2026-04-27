DROP INDEX IF EXISTS idx_expense_tag_tag_id;
DROP TABLE IF EXISTS expense_tag;

DROP INDEX IF EXISTS idx_expense_deleted_at;
DROP INDEX IF EXISTS idx_expense_expense_date;
DROP INDEX IF EXISTS idx_expense_payment_method_id;
DROP INDEX IF EXISTS idx_expense_category_id;
DROP TABLE IF EXISTS expense;

DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS payment_method;
DROP TABLE IF EXISTS category;

DROP TYPE IF EXISTS payment_method_type;
