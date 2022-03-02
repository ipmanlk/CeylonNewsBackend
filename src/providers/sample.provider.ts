import axios from "axios";
import cheerio from "cheerio";
import { Source, Post } from "../types/main";
import { extract } from "article-parser";
import { getPrismaClient } from "../services/database.service.js";
import { writeFileSync } from "fs";

const prisma = getPrismaClient();

export async function scrape(source: Source) {
	try {
		const { data: pageSource } = await axios.get(source.url);

		const $page = cheerio.load(pageSource);
		const postUrls: Set<string> = new Set();

		$page(".all-section-tittle").each((_, el) => {
			const url = $page(el).children("a").first().attr("href");
			if (!url || url.trim() === "") return;
			postUrls.add(encodeURI(url));
			return postUrls.size < 10 ? true : false;
		});

		const posts: Post[] = [];

		for (const url of postUrls) {
			const dbPost = await prisma.post.findFirst({
				select: { url: true },
				where: { url: url },
			});

			if (!process.env.DEVELOPMENT && dbPost) continue;

			const article = await extract(url).catch(console.log);

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
		console.log(e);
		return [];
	}
}

scrape({
	name: "Neth News",
	language: "si",
	url: "https://www.hirunews.lk/local-news.php?pageID=1",
	provider: "rss_provider.ts",
	enabled: true,
}).then((posts) => {
	writeFileSync("test.json", JSON.stringify(posts));
});
