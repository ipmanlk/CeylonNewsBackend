import axios from "axios";
import { Source, Post } from "../types/main";
import { extract } from "article-parser";
import { getPrismaClient } from "../services/database.service.js";
import { XMLParser } from "fast-xml-parser";

const prisma = getPrismaClient();

export async function scrape(source: Source) {
	try {
		const { data: pageSource } = await axios.get(source.url);

		const parser = new XMLParser();
		const obj = parser.parse(pageSource);

		const postUrls: Set<string> = new Set();
		obj.rss.channel.item
			.slice(0, 5)
			.forEach((u: any) => postUrls.add(encodeURI(decodeURI(u.link))));

		const posts: Post[] = [];

		for (const url of postUrls) {
			const dbPost = await prisma.post.findFirst({
				select: { url: true },
				where: { url: url },
			});

			if (!process.env.DEVELOPMENT && dbPost) continue;

			const article = await extract(url).catch((_) => {});
			if (!article || !article.title || !article.content) continue;

			const post: Post = {
				title: article.title,
				content: article.content,
				url: url,
				thumbnailUrl: article.image,
				language: source.language,
				sourceName: source.name,
				createdDate:
					article.published && article.published.trim() !== ""
						? new Date(article.published)
						: new Date(),
			};
			posts.push(post);
		}

		return posts;
	} catch (e) {
		return [];
	}
}
