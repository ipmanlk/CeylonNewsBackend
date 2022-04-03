import { assert } from "chai";
import { Provider } from "../src/types/main";
import { sources } from "../src/sources.js";

describe("news sources", () => {
	for (const source of sources) {
		it(`source: ${source.name}-${source.language} | provider: ${source.provider} | url: ${source.url}`, async () => {
			const provider: Provider = await import(
				`../src/providers/${source.provider}.ts`
			);
			const posts = await provider.scrape(source);

			assert.notEqual(posts.length, 0);
		}).timeout(80000);
	}
});
