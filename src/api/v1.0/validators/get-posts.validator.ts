import { Request } from "express";
import { sources } from "../../../sources";
import { Language } from "../../../types/main";

export function validate(req: Request) {
	const query = req.query;

	let keyword = "";

	if (query.keyword && typeof query.keyword === "string") {
		keyword = query.keyword;
	}

	if (!query.languages || typeof query.languages !== "string") {
		return {
			error: "Please provide a valid list of languages.",
		};
	}

	const inputLanguages = query.languages
		.split(",")
		.map((i) => i.toLowerCase().trim());

	let hasInvalidLanguage = false;

	inputLanguages.every((lang) => {
		if (!["si", "en", "ta"].includes(lang)) {
			hasInvalidLanguage = true;
			return false;
		}
		return true;
	});

	if (hasInvalidLanguage) {
		return {
			error: "Your request contains an invalid language.",
		};
	}

	let inputSources: string[] = [];

	if (query.sources) {
		if (typeof query.sources !== "string") {
			return {
				error: "Your request contains an invalid sources.",
			};
		}

		inputSources = query.sources.split(",").map((i) => i.toLowerCase().trim());

		let hasInvalidSource = false;
		let invalidSourceName = "";

		inputSources.every((inputSource) => {
			const result = sources.find((s) => s.name.toLowerCase() == inputSource);

			if (!result) {
				hasInvalidSource = true;
				invalidSourceName = inputSource;
				return false;
			}
			return true;
		});

		if (hasInvalidSource) {
			return {
				error: `Your request contains an invalid source ${invalidSourceName}.`,
			};
		}
	}

	let skip = 0;
	let limit = 20;

	if (
		(query.limit && typeof query.limit === "string") ||
		query.limit === "number"
	) {
		const inputLimit = parseInt(query.limit);

		if (isNaN(inputLimit)) {
			return {
				error: `Please provide a valid number for the limit.`,
			};
		}

		if (inputLimit <= 15) limit = inputLimit;
	}

	if (
		(query.skip && typeof query.skip === "string") ||
		query.skip === "number"
	) {
		const inputSkip = parseInt(query.skip);

		if (isNaN(inputSkip)) {
			return {
				error: `Please provide a valid number for the skip.`,
			};
		}

		skip = inputSkip;
	}

	return {
		data: {
			keyword: keyword,
			sources: inputSources,
			languages: inputLanguages as Array<Language>,
			limit: limit,
			skip: skip,
		},
	};
}
