import { Request, Response } from "express";
import { sources } from "../../../sources";
import { validate } from "../validators/get-sources.validator";

export function getSources(req: Request, res: Response) {
	const validatedInputs = validate(req);

	if (validatedInputs.error || !validatedInputs.data) {
		return res.status(400).json(validatedInputs);
	}

	res.status(200).json(
		sources
			.filter((s) => validatedInputs.data.languages.includes(s.language))
			.map((s) => {
				return {
					name: s.name,
				};
			})
	);
}
