export type Language = "si" | "en" | "ta"; //ISO 639-1

export interface Source {
	name: string;
	language: Language;
	url: string;
	provider: string;
	enabled: boolean;
}

export interface Post {
	title: string;
	content: string;
	url: string;
	thumbnailUrl?: string;
	language: Language;
	sourceName: string;
	createdDate: Date;
}

export interface Provider {
	scrape: (source: Source) => Array<Post>;
}
