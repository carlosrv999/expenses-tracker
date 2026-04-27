export function formatAmount(amountMinor: number, currency = 'PEN'): string {
	const major = amountMinor / 100;
	try {
		return new Intl.NumberFormat('es-PE', {
			style: 'currency',
			currency
		}).format(major);
	} catch {
		return `${currency} ${major.toFixed(2)}`;
	}
}

export function toMinorAmount(value: string | number): number {
	const n = typeof value === 'number' ? value : Number(value);
	if (!Number.isFinite(n)) return 0;
	return Math.round(n * 100);
}

export function toMajorAmount(minor: number): number {
	return minor / 100;
}

export function formatDate(iso: string): string {
	if (!iso) return '';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return iso;
	return d.toLocaleString('es-PE', {
		year: 'numeric',
		month: 'short',
		day: '2-digit',
		hour: '2-digit',
		minute: '2-digit'
	});
}

export function toDateTimeLocal(iso: string | null | undefined): string {
	if (!iso) {
		const now = new Date();
		now.setSeconds(0, 0);
		return toLocalInputValue(now);
	}
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '';
	return toLocalInputValue(d);
}

function toLocalInputValue(d: Date): string {
	const pad = (n: number) => String(n).padStart(2, '0');
	return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function fromDateTimeLocal(value: string): string {
	if (!value) return '';
	const d = new Date(value);
	return d.toISOString();
}
