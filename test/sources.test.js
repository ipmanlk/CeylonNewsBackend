const Parser = require("rss-parser");
const parser = new Parser();
const sources = require("../sources/sources.json");

// settings for jest
jest.setTimeout(120000);

// Test all RSS feeds
for (let source of sources) {
    test(`RSS Feed: ${source.name} (${source.lang})`, async() => {
        console.log(`Testing ${source.name} (${source.lang})`);
        const feed = await parser.parseURL(source.feed);
        expect(feed.items.length).not.toBe(0);
    });
}