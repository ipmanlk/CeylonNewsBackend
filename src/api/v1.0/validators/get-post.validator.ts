import { Request } from "express";

export function validate(req: Request) {
	const params = req.params;

	if (
		!params.id ||
		(typeof params.id !== "string" && typeof params.id !== "number")
	) {
		return { error: "Please provide a valid news id." };
	}

	return {
		data: {
			id: parseInt(params.id),
		},
	};
}
