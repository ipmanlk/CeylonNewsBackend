const { SQLite } = require("../../libs/dmxSQLite");
const db = new SQLite(`${__dirname}/../../db/cn.db`);

const getNewsList = (lang, sources = [], lastNewsId = false) => {
    if (sources.length !== 0) {
        let sourcesStr = "'" + sources.toString().replace(/\,/g, "','") + "'";
        let sql;

        if (lastNewsId) {
            // for load more option
            sql = `SELECT n.id, n.title, n.time, s.name as source FROM news n,source s WHERE n.id < '${lastNewsId}' AND s.lang = '${lang}' AND n.source_id IN (${sourcesStr}) AND s.id = n.source_id ORDER BY n.id DESC LIMIT 8;`;
        } else {
            // initial news list
            sql = `SELECT n.id, n.title, n.time, s.name as source FROM news n,source s WHERE s.lang = '${lang}' AND n.source_id IN (${sourcesStr}) AND s.id = n.source_id ORDER BY n.id DESC LIMIT 8;`;
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
    let sql = `SELECT id, name, desc FROM source WHERE lang = '${lang}'`;        
    return db.getAll(sql);
}

module.exports = {
    getNewsList,
    getNewsPost,
    getNewsSources
}