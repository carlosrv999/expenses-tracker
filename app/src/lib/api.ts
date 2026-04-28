import type {
	Category,
	CategoryInput,
	Expense,
	ExpenseInput,
	ExpenseListFilters,
	PaymentMethod,
	PaymentMethodInput,
	Tag,
	TagInput
} from './types';

const BASE_URL = 'http://localhost:8080';

// Add the paginated interface to match the backend PaginatedExpenseList struct
export interface PaginatedExpenseList {
	expenses: Expense[];
	total_count: number;
	limit: number;
	offset: number;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`${BASE_URL}${path}`, {
		...init,
		headers: {
			'Content-Type': 'application/json',
			...(init?.headers ?? {})
		}
	});
	if (!res.ok) {
		let detail = '';
		try {
			detail = await res.text();
		} catch {
			// ignore
		}
		throw new Error(`HTTP ${res.status}: ${detail || res.statusText}`);
	}
	if (res.status === 204) return undefined as T;
	const text = await res.text();
	if (!text) return undefined as T;
	return JSON.parse(text) as T;
}

function toQuery(params: Record<string, unknown>): string {
	const usp = new URLSearchParams();
	for (const [key, value] of Object.entries(params)) {
		if (value === undefined || value === null || value === '') continue;
		usp.set(key, String(value));
	}
	const s = usp.toString();
	return s ? `?${s}` : '';
}

export const categoriesApi = {
	list: () => request<Category[]>('/api/v1/categories'),
	get: (id: number) => request<Category>(`/api/v1/categories/${id}`),
	create: (body: CategoryInput) =>
		request<Category>('/api/v1/categories', {
			method: 'POST',
			body: JSON.stringify(body)
		}),
	update: (id: number, body: CategoryInput) =>
		request<Category>(`/api/v1/categories/${id}`, {
			method: 'PUT',
			body: JSON.stringify(body)
		}),
	remove: (id: number) =>
		request<void>(`/api/v1/categories/${id}`, {
			method: 'DELETE'
		})
};

export const tagsApi = {
	list: () => request<Tag[]>('/api/v1/tags'),
	get: (id: number) => request<Tag>(`/api/v1/tags/${id}`),
	create: (body: TagInput) =>
		request<Tag>('/api/v1/tags', {
			method: 'POST',
			body: JSON.stringify(body)
		}),
	update: (id: number, body: TagInput) =>
		request<Tag>(`/api/v1/tags/${id}`, {
			method: 'PUT',
			body: JSON.stringify(body)
		}),
	remove: (id: number) =>
		request<void>(`/api/v1/tags/${id}`, {
			method: 'DELETE'
		})
};

export const paymentMethodsApi = {
	list: () => request<PaymentMethod[]>('/api/v1/payment-methods'),
	get: (id: number) => request<PaymentMethod>(`/api/v1/payment-methods/${id}`),
	create: (body: PaymentMethodInput) =>
		request<PaymentMethod>('/api/v1/payment-methods', {
			method: 'POST',
			body: JSON.stringify(body)
		}),
	update: (id: number, body: PaymentMethodInput) =>
		request<PaymentMethod>(`/api/v1/payment-methods/${id}`, {
			method: 'PUT',
			body: JSON.stringify(body)
		}),
	remove: (id: number) =>
		request<void>(`/api/v1/payment-methods/${id}`, {
			method: 'DELETE'
		})
};

export const expensesApi = {
	// Updated return type to PaginatedExpenseList
	list: (filters: ExpenseListFilters = {}) =>
		request<PaginatedExpenseList>(`/api/v1/expenses${toQuery(filters as Record<string, unknown>)}`),
	get: (id: number) => request<Expense>(`/api/v1/expenses/${id}`),
	create: (body: ExpenseInput) =>
		request<Expense>('/api/v1/expenses', {
			method: 'POST',
			body: JSON.stringify(body)
		}),
	update: (id: number, body: ExpenseInput) =>
		request<Expense>(`/api/v1/expenses/${id}`, {
			method: 'PUT',
			body: JSON.stringify(body)
		}),
	remove: (id: number) =>
		request<void>(`/api/v1/expenses/${id}`, {
			method: 'DELETE'
		})
};
