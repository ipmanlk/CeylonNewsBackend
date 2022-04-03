import { getPrismaClient } from "./services/database.service.js";
import { sources } from "./sources.js";
import { Post, Provider } from "./types/main";

const prisma = getPrismaClient();

export async function startScraper() {
	await scrapeSources();
	setInterval(scrapeSources, 15 * 60000);
}

async function scrapeSources() {
	let posts: Array<Post> = [];

	for (const source of sources) {
		try {
			console.log(`scraping: ${source.name}-${source.language}`);
			const provider: Provider = await import(`./providers/${source.provider}.js`);
			const sourcePosts = await provider.scrape(source);
			posts = posts.concat(sourcePosts);
			await Promise.all(
				sourcePosts.map((p) =>
					prisma.post.create({
						data: {
							title: p.title,
							content: p.content,
							url: p.url,
							thumbnailUrl: p.thumbnailUrl,
							language: p.language,
							sourceName: p.sourceName,
							createdDate: p.createdDate,
						},
					})
				)
			);
		} catch (e) {
			console.log(e);
		}
	}
}
