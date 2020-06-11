const Parser = require("rss-parser");
const parser = new Parser();
const sources = require("../src/sources/sources.json");

// settings for jest
jest.setTimeout(180000);

// Test all RSS feeds
for (let source of sources) {
    if (!source.enabled) continue;

    if (source.site) {
        test(`RSS Feed: ${source.name} (${source.lang})`, async () => {
            console.log(`Testing scraper ${source.name} (${source.lang})`);
            const scraper = require(`${__dirname}/../src/scraper/sites/${source.scraper}.js`);
            const posts = await scraper.scrape(source);
            expect(posts.length).not.toBe(0);
        });
        continue;
    }
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
        return source.feed || source.site;
    });

    const duplicates = hasDuplicates(namesAndLangs) || hasDuplicates(feeds);

    expect(duplicates).toBe(false);
});

const hasDuplicates = (array) => {
    return (new Set(array)).size !== array.length;
}