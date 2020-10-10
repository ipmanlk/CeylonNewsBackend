const { SQLite } = require("../libs/dmxSQLite");
const cheerio = require('cheerio');
const CronJob = require('cron').CronJob;
const Parser = require('rss-parser');
const db = new SQLite(`${__dirname}/../../db/cn.db`);
const parser = new Parser();
const storeNews = require("./storeNews");
const sources = require("../sources/sources.json");

const start = async () => {
	try {
		scrapeNews();
	} catch (error) {
		throw Error(error);
	}
}

// remove news from database and insert sources to db from sources.json
const prepareDB = async () => {
	await db.run("DELETE FROM news");
	await db.run("DELETE FROM source");
	await db.run("UPDATE sqlite_sequence SET seq = 0");

	for (let source of sources) {
		await db.run(`INSERT INTO source(id, name, desc, lang, enabled) VALUES (?,?,?,?,?)`, [source.id, source.name, source.desc, source.lang, source.enabled]);
	}
}

const scrapeNews = async () => {
	// check news table count
	storeNews.checkNewsCount();

	const sites = sources.filter(source => source.site);
	const rssFeeds = sources.filter(source => source.feed);

	await scrapeFeeds(rssFeeds);
	await scrapeSites(sites);
}

const scrapeFeeds = async (feeds) => {
	const articles = [];

	for (let feedData of feeds) {
		const feed = await parser.parseURL(feedData.feed);

		// loop through posts in each rss feed
		feed.items.forEach(item => {
			const main_img = getMainImg(item.content);

			const article = {
				source_id: feedData.id,
				title: item.title,
				post: item.content,
				link: item.link,
				main_img: main_img
			};

			articles.push(clean(article));
		});
	}

	storeNews.save(articles);
}

const scrapeSites = async (sites) => {
	for (let siteData of sites) {
		const scraper = require(`${__dirname}/sites/${siteData.scraper}.js`);
		const articles = await scraper.scrape(siteData);
		storeNews.save(articles);
	}
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

const scrapeCronJob = new CronJob("*/10 * * * *", () => {
	try {
		start();
	} catch (e) {
		console.log(e);
	}
}, null, true, 'Asia/Colombo');

// remove useless elements from news articles
const clean = (article) => {
	const $ = cheerio.load(article.post);
	try {
		$('.sr-date').parent().remove();
		$('a[href="https://blockads.fivefilters.org"]').parent().remove();
		$('a[href="https://blockads.fivefilters.org/acceptable.html"]').parent().remove();
		$('img[src="https://data.gossiplankanews.com/box0/arti.png"]').remove();
		$(".adsbygoogle").remove();
		$('*').removeAttr("style");
		article.post = escapedHtmlFix($.html());
		article.title = escapedHtmlFix(article.title);
	} catch (e) {
		console.log(e);
	}

	return article;
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
}

module.exports = {
	scrapeCronJob,
	prepareDB
}