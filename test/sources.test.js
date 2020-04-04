const Parser = require("rss-parser");
const parser = new Parser();
const sources = require("../sources/sources.json");

// settings for jest
jest.setTimeout(180000);

// Test all RSS feeds
for (let source of sources) {
    test(`RSS Feed: ${source.name} (${source.lang})`, async () => {
        console.log(`Testing ${source.name} (${source.lang})`);
        const feed = await parser.parseURL(source.feed);
        expect(feed.items.length).not.toBe(0);
    });
}

// Test for duplicate entries in sources.json
test("check for duplicate sources", () => {
    const namesAndLangs = sources.map(source => {
        return source.name + source.lang;
    });

    const feeds = sources.map(source => {
        return source.feed;
    });

    const duplicates = hasDuplicates(namesAndLangs) || hasDuplicates(feeds);

    expect(duplicates).toBe(false);
});

const hasDuplicates = (array) => {
    return (new Set(array)).size !== array.length;
}