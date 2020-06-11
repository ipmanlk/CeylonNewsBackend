const { SQLite } = require("../libs/dmxSQLite");
const cheerio = require('cheerio');
const moment = require('moment-timezone');
const db = new SQLite(`${__dirname}/../../db/cn.db`);

const save = async (newsData) => {
    let news = clean(newsData);
    if (!news.title || news.title == "" || news.title.trim() == "") return;
    const currentDateTime = moment().tz('Asia/Colombo').format('YYYY-MM-DD hh:mm A');

    // check if news item is already in db
    const row = await db.getOne("SELECT COUNT(id) AS count FROM news WHERE title = ? OR link = ?", [news.title, news.link]);
    if (row.count !== 0) return;

    // else save it 
    await db.run("INSERT INTO news(title, post, link, source_id, time, main_img) VALUES(?,?,?,?,?,?)", [news.title, news.post, news.link, news.source_id, currentDateTime, news.main_img]);
}

// remove useless elements from news item
const clean = (newsData) => {
    const $ = cheerio.load(newsData.post);
    try {
        $('.sr-date').parent().remove();
        $('a[href="https://blockads.fivefilters.org"]').parent().remove();
        $('a[href="https://blockads.fivefilters.org/acceptable.html"]').parent().remove();
        $('img[src="https://data.gossiplankanews.com/box0/arti.png"]').remove();
        $(".adsbygoogle").remove();
        $('*').removeAttr("style");
        newsData.post = escapedHtmlFix($.html());
        newsData.title = escapedHtmlFix(newsData.title);
    } catch (e) {
        console.log(e);
    }

    return newsData;
}

// fix broken/escaped html tags and elements
const escapedHtmlFix = (text) => {
    return text
        .replace("&amp;", "&")
        .replace("&lt;", "<")
        .replace("&gt;", ">")
        .replace("&quot;", '"')
        .replace("&#039;", "'")
        .replace("&amp;#039;", "'");
};

const checkNewsCount = async () => {
    const row = await db.getOne("SELECT COUNT(id) as count FROM news");
    // empty news table if there are too many news
    if (row.count > 10000) {
        await db.run("DELETE FROM news");
    }
}

module.exports = {
    save,
    checkNewsCount
}