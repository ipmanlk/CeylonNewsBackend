import axios from "axios";
import cheerio from "cheerio";
import { Source, Post } from "../types/main";
import { extract } from "article-parser";
import { getPrismaClient } from "../services/database.service.js";

const prisma = getPrismaClient();

export async function scrape(source: Source) {
	try {
		const { data: pageSource } = await axios.get(source.url);

		const $page = cheerio.load(pageSource);
		const postUrls: Set<string> = new Set();

		$page(".itemList .itemContainer .catItemTitle a").each((_, el) => {
			const path = $page(el).attr("href");
			if (!path || path.trim() === "") return;

			let url = "";
			if (source.language === "en") {
				url = "https://news.lk" + path;
			} else if (source.language === "ta") {
				url = "https://tamil.news.lk" + path;
			} else {
				url = "https://sinhala.news.lk" + path;
			}

			postUrls.add(encodeURI(decodeURI(url)));
			return postUrls.size < 5 ? true : false;
		});

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
		console.error(e);
		return [];
	}
}
