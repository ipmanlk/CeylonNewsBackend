const { SQLite } = require("../../libs/dmxSQLite");
const dbPath = `${__dirname}/../../../db/cn.db`
const db = new SQLite(dbPath);

const getNewsList = (sources = [], lastNewsId = false) => {
    if (sources.length !== 0) {
        let sourcesStr = "'" + sources.toString().replace(/\,/g, "','") + "'";
        let sql;

        if (lastNewsId) {
            // for load more option
            sql = `SELECT n.id, n.title, n.time, n.main_img, s.name as source FROM news n,source s WHERE n.id < '${lastNewsId}' AND n.source_id IN (${sourcesStr}) AND s.id = n.source_id ORDER BY n.id DESC LIMIT 8;`;
        } else {
            // initial news list
            sql = `SELECT n.id, n.title, n.time, n.main_img, s.name as source FROM news n,source s WHERE n.source_id IN (${sourcesStr}) AND s.id = n.source_id ORDER BY n.id DESC LIMIT 20;`;
        }

        return db.getAll(sql);
    } else {
        return [];
    }
}

const getNewsPost = (newsId) => {
    let sql = `SELECT post, link FROM news WHERE id = ${newsId}`;
    return db.getOne(sql);
}

const getNewsSources = (lang) => {
    let sql = `SELECT id, name, desc, enabled FROM source WHERE lang = '${lang}'`;
    return db.getAll(sql);
}

const getLatestId = (sources = []) => {
    let sourcesStr = "'" + sources.toString().replace(/\,/g, "','") + "'";
    let sql = `SELECT n.id FROM news n,source s WHERE n.source_id IN (${sourcesStr}) AND s.id = n.source_id ORDER BY n.id DESC LIMIT 1;`;
    return db.getAll(sql);
}

module.exports = {
    getNewsList,
    getNewsPost,
    getNewsSources,
    getLatestId
}