const { SQLite } = require("../../../libs/dmxSQLite");
const dbPath = `${__dirname}/../../../../db/cn.db`
const db = new SQLite(dbPath);

const getNewsList = (sources = [], keyword = "", skip = 0) => {
    // convert sources array to a string
    const sourcesStr = "'" + sources.toString().replace(/\,/g, "','") + "'";

    // prepare query to search db
    const sql = `SELECT n.id, n.title, n.time, n.main_img, s.name as source FROM news n,source s WHERE n.source_id IN (${sourcesStr}) AND s.id = n.source_id AND n.title LIKE '%${keyword}%' OR n.time LIKE '${keyword}' ORDER BY n.id DESC LIMIT ${skip}, 20;`

    return db.getAll(sql);
}

const getNewsPost = (newsId) => {
    let sql = `SELECT n.title, n.time, n.main_img, n.post, s.name as source_name, n.link FROM news n, source s WHERE n.id = ${newsId} AND n.source_id = s.id`;
    return db.getOne(sql);
}

const getNewsSources = (lang) => {
    let sql = `SELECT id, name, desc, enabled FROM source WHERE lang = '${lang}'`;
    return db.getAll(sql);
}

module.exports = {
    getNewsList,
    getNewsSources,
    getNewsPost
}