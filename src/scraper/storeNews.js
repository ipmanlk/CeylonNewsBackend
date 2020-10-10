const moment = require("moment-timezone");
const sqlite3 = require("sqlite3").verbose();
let db = new sqlite3.Database(`${__dirname}/../../db/cn.db`, sqlite3.OPEN_READWRITE);

const save = (articles) => {
	const currentDateTime = moment().tz('Asia/Colombo').format('YYYY-MM-DD hh:mm A');

	db.serialize(function () {
		db.run("BEGIN TRANSACTION");

		for (let article of articles) {
			if (!article.title || article.title.trim() == "" || article.title.length < 10) continue;
			db.run("INSERT OR IGNORE INTO news(title, post, link, source_id, time, main_img) VALUES (?,?,?,?,?,?)", [
				article.title,
				article.post,
				article.link,
				article.source_id,
				currentDateTime,
				article.main_img
			]);
		}

		db.run("COMMIT");
	});
}

const checkNewsCount = () => {
	db.get("SELECT COUNT(id) as count FROM news", [], (err, row) => {
		if (row.count > 100000) {
			db.run("DELETE FROM news");
		}
	});
}

module.exports = {
	save,
	checkNewsCount
}