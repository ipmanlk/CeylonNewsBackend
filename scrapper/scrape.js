const { SQLite } = require("../libs/dmxSQLite");
const Parser = require('rss-parser');
const db = new SQLite(`${__dirname}/../db/cn.db`);
const parser = new Parser();
const storeNews = require("./storeNews");

const start = async () => {
    scrapeNews().catch(e => console.log(e));
}

const scrapeNews = async () => {
    const sources = await getSources().catch(e => console.log(e));

    // loop through each rss feed
    sources.forEach(async(source) => {
        const feed = await parser.parseURL(source.feed);

        // loop through posts in each rss feed
        feed.items.forEach(item => {
            let newsData = {
                source_id: source.id,
                title: item.title,
                post: item.content,
                link: item.link
            };

            storeNews.save(newsData).catch(e => {
                console.log(e);
            });

        });
    });
}

const getSources = async () => {
    const sources = await db.getAll("SELECT * FROM source").catch(e => console.log(e));
    return sources;
}

start().catch(e => {
    console.log(e);
});