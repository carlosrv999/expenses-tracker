export type PaymentMethodType =
	| 'credit_card'
	| 'cash'
	| 'debit_card'
	| 'yape'
	| 'plin'
	| 'bank_transfer';

export interface Category {
	category_id: number;
	parent_category_id: number | null;
	category_name: string;
	icon: string | null;
	color: string | null;
	created_at: string;
	updated_at: string;
}

export interface CategoryInput {
	parent_category_id?: number | null;
	category_name: string;
	icon?: string | null;
	color?: string | null;
}

export interface Tag {
	tag_id: number;
	tag_name: string;
	color: string | null;
	icon: string | null;
	created_at: string;
	updated_at: string;
}

export interface TagInput {
	tag_name: string;
	color?: string | null;
	icon?: string | null;
}

export interface PaymentMethod {
	payment_method_id: number;
	method_name: string;
	method_type: PaymentMethodType;
	icon: string | null;
	created_at: string;
	updated_at: string;
}

export interface PaymentMethodInput {
	method_name: string;
	method_type: PaymentMethodType;
	icon?: string | null;
}

export interface Expense {
	expense_id: number;
	category_id: number;
	payment_method_id: number;
	currency: string;
	amount: number;
	expense_date: string;
	merchant_name: string | null;
	description: string | null;
	created_at: string;
	updated_at: string;
	deleted_at: string | null;
	tag_ids?: number[];
	tags?: Tag[];
}

export interface ExpenseInput {
	category_id: number;
	payment_method_id: number;
	currency: string;
	amount: number;
	expense_date: string;
	merchant_name?: string | null;
	description?: string | null;
	tag_ids?: number[];
}

export interface ExpenseListFilters {
	category_id?: number;
	payment_method_id?: number;
	start_date?: string;
	end_date?: string;
	tags?: number[];
	limit?: number;
	offset?: number;
}
