import { Request } from "express";

export function validate(req: Request) {
	const query = req.query;

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

	return {
		data: {
			languages: inputLanguages,
		},
	};
}
