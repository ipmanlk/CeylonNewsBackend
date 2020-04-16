const { SQLite } = require("../libs/dmxSQLite");
const CronJob = require('cron').CronJob;
const Parser = require('rss-parser');
const db = new SQLite(`${__dirname}/../../db/cn.db`);
const parser = new Parser();
const storeNews = require("./storeNews");

const start = async () => {
    try {
        scrapeNews();
    } catch (error) {
        throw Error(error);
    }
}

// remove news from database and insert sources to db from sources.json
const prepareDB = async() => {
    await db.run("DELETE FROM news");
    await db.run("DELETE FROM source");
    await db.run("UPDATE sqlite_sequence SET seq = 0");

    const sources = require("../sources/sources.json");
    sources.forEach(async(source) => {
        await db.run(`INSERT INTO source(id, name, desc, lang, feed, enabled) VALUES (?,?,?,?,?,?)`, [source.id, source.name, source.desc, source.lang, source.feed, source.enabled]);
    });
}

const scrapeNews = async () => {
    // check news table count
    await storeNews.checkNewsCount();

    const sources = await getSources();

    // loop through each rss feed
    sources.forEach(async (source) => {
        const feed = await parser.parseURL(source.feed);

        // loop through posts in each rss feed
        feed.items.forEach(item => {
            const main_img = getMainImg(item.content);

            const newsData = {
                source_id: source.id,
                title: item.title,
                post: item.content,
                link: item.link,
                main_img: main_img
            };

            storeNews.save(newsData).catch(e => {});
        });
    });
}

const getMainImg = (post) => {
    const regEx = /<img.+src\=(?:\"|\')(.+?)(?:\"|\')(?:.+?)\>/;
    try {
        let imgs = (regEx.exec(`${post}`));
        let img = imgs[3] || imgs[2] || imgs[1];
        img = img.replace(/^http:\/\//i, 'https://');
        return img;
    } catch (error) {
        return "null";
    }
}

const getSources = async () => {
    const sources = await db.getAll("SELECT * FROM source");
    return sources;
}

const scrapeCronJob = new CronJob("*/1 * * * *", () => {
    try {
        start();
    } catch (e) {
        console.log(e);
    }
}, null, true, 'Asia/Colombo');

module.exports = {
    scrapeCronJob,
    prepareDB
}