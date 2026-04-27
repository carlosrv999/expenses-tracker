import type { Component } from 'svelte';
import {
	Wallet,
	Banknote,
	CreditCard,
	Smartphone,
	Phone,
	Landmark,
	ArrowLeftRight,
	Utensils,
	UtensilsCrossed,
	Coffee,
	ShoppingCart,
	ShoppingBag,
	Car,
	Bus,
	Plane,
	Home,
	Lightbulb,
	Briefcase,
	BriefcaseBusiness,
	Heart,
	Stethoscope,
	GraduationCap,
	Gamepad2,
	Film,
	Music,
	Gift,
	Dog,
	Baby,
	Dumbbell,
	Fuel,
	Tag as TagIcon,
	Folder,
	FolderTree,
	Tags,
	Receipt
} from 'lucide-svelte';

// lucide-svelte ships legacy Svelte types; coerce to Svelte 5's Component type for runes mode.
export type IconComponent = Component<{ size?: number | string; color?: string; strokeWidth?: number | string }>;

export const ICON_MAP: Record<string, IconComponent> = {
	wallet: Wallet as unknown as IconComponent,
	'money-bill': Banknote as unknown as IconComponent,
	money: Banknote as unknown as IconComponent,
	cash: Banknote as unknown as IconComponent,
	'credit-card': CreditCard as unknown as IconComponent,
	mobile: Smartphone as unknown as IconComponent,
	smartphone: Smartphone as unknown as IconComponent,
	phone: Phone as unknown as IconComponent,
	bank: Landmark as unknown as IconComponent,
	'bank-transfer': ArrowLeftRight as unknown as IconComponent,
	transfer: ArrowLeftRight as unknown as IconComponent,
	utensils: Utensils as unknown as IconComponent,
	restaurant: UtensilsCrossed as unknown as IconComponent,
	coffee: Coffee as unknown as IconComponent,
	'shopping-cart': ShoppingCart as unknown as IconComponent,
	'shopping-bag': ShoppingBag as unknown as IconComponent,
	car: Car as unknown as IconComponent,
	bus: Bus as unknown as IconComponent,
	plane: Plane as unknown as IconComponent,
	home: Home as unknown as IconComponent,
	utilities: Lightbulb as unknown as IconComponent,
	briefcase: Briefcase as unknown as IconComponent,
	'briefcase-business': BriefcaseBusiness as unknown as IconComponent,
	health: Heart as unknown as IconComponent,
	medical: Stethoscope as unknown as IconComponent,
	education: GraduationCap as unknown as IconComponent,
	gaming: Gamepad2 as unknown as IconComponent,
	movie: Film as unknown as IconComponent,
	music: Music as unknown as IconComponent,
	gift: Gift as unknown as IconComponent,
	pet: Dog as unknown as IconComponent,
	baby: Baby as unknown as IconComponent,
	gym: Dumbbell as unknown as IconComponent,
	fuel: Fuel as unknown as IconComponent,
	tag: TagIcon as unknown as IconComponent,
	tags: Tags as unknown as IconComponent,
	folder: Folder as unknown as IconComponent,
	'folder-tree': FolderTree as unknown as IconComponent,
	receipt: Receipt as unknown as IconComponent
};

export const AVAILABLE_ICONS = Object.keys(ICON_MAP).sort();

export function resolveIcon(
	name: string | null | undefined,
	fallback: IconComponent
): IconComponent {
	if (!name) return fallback;
	return ICON_MAP[name] ?? fallback;
}
