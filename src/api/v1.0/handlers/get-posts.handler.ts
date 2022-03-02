import { Post } from "@prisma/client";
import { Request, Response } from "express";
import { getPrismaClient } from "../../../services/database.service";
import { validate } from "../validators/get-posts.validator";

const prisma = getPrismaClient();

export async function getPosts(req: Request, res: Response) {
	const validatedInputs = validate(req);

	if (validatedInputs.error || !validatedInputs.data) {
		return res.status(400).json(validatedInputs);
	}

	try {
		let posts: Post[] = [];

		if (validatedInputs.data.sources.length > 0) {
			posts = await prisma.post.findMany({
				where: {
					OR: {
						title: {
							contains: validatedInputs.data.keyword,
						},
						content: {
							contains: validatedInputs.data.keyword,
						},
					},
					sourceName: {
						in: validatedInputs.data.sources,
					},
					language: {
						in: validatedInputs.data.languages,
					},
				},
				take: validatedInputs.data.limit,
				skip: validatedInputs.data.skip,
				orderBy: { createdDate: "desc" },
			});
		} else {
			posts = await prisma.post.findMany({
				where: {
					OR: {
						title: {
							contains: validatedInputs.data.keyword,
						},
						content: {
							contains: validatedInputs.data.keyword,
						},
					},
					language: {
						in: validatedInputs.data.languages,
					},
				},
				take: validatedInputs.data.limit,
				skip: validatedInputs.data.skip,
				orderBy: { createdDate: "desc" },
			});
		}

		return res.status(200).json(posts);
	} catch (e) {
		res.status(500).json({ error: "Something went wrong." });
	}
}
