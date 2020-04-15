const tokens = require("../auth/v2.0/tokens.json");
const newsData = require("./data/newsData");
const { stLogger } = require("sematext-agent-express");

const handle = (req, res) => {

    // check authentication
    const token = req.headers.token;
    if (tokens.indexOf(token) == -1) {
        stLogger.error("fail: authentication");
        res.status(401).json({"error": "You are not authorized."});
        return;
    }

    const action = req.query.action;

    // search & navigate news list 
    if (action == "news-list") {
        const sources = req.query.sources.split(",") || [];
        const keyword = req.query.keyword || "";
        const skip = req.query.skip || 0;
        newsData.getNewsList(sources, keyword, skip).then(data => {
            stLogger.info("sent: news-list")
            res.json(data);
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-list server error");
            res.status(500).json({ "error": "Internal Server Error." });
        });
        return;
    }

    // retrieve single news post
    if (action == "news-post") {
        newsData.getNewsPost(req.query.news_id).then(newsPost => {
            if (newsPost == undefined) {
                stLogger.error("fail: news-post not found with given id")
                res.status(404).json({ "error": "Unable to locate that post." });
            } else {
                stLogger.info("sent: news post")
                res.json(newsPost);
            }
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-post server error");
            res.status(500).json({ "error": "Internal Server Error." });
        });
        return;
    }

    // retrieve news sources list
    if (action == "news-sources" && req.query.lang) {
        newsData.getNewsSources(req.query.lang).then(newsSources => {
            stLogger.info("sent: news-sources");
            res.json(newsSources);
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-sources server error")
            res.status(500).json({ "error": "Internal Server Error." });
        });
        return;
    }

    // default response
    res.status(404).json({ "error": "Sorry can't find that!" });

}

module.exports = {
    handle
}