import { Request, Response } from "express";
import { getPrismaClient } from "../../../services/database.service";
import { validate } from "../validators/get-post.validator";

const prisma = getPrismaClient();

export async function getPost(req: Request, res: Response) {
	const validatedInputs = validate(req);

	if (validatedInputs.error || !validatedInputs.data) {
		return res.status(400).json(validatedInputs);
	}

	try {
		const post = await prisma.post.findFirst({
			where: { id: validatedInputs.data.id },
		});

		if (post) {
			res.status(200).json(post);
		} else {
			res
				.status(404)
				.json({ error: "Unable to find a news article with the given id." });
		}
	} catch (e) {
		res.status(500).json({ error: "Something went wrong." });
	}
}
