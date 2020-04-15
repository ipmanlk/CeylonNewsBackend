const { stLogger } = require("sematext-agent-express");
const getNews = require("./getNews");

const handle = (req, res) => {
    const params = req.query;

    // news post
    if (params.action == "news-post" && params.post_id) {
        getNews.getNewsPost(params.post_id).then(newsPost => {
            if (newsPost == undefined) {
                stLogger.error("fail: news-post not found with given id");
                res.status(404).send(JSON.stringify({ "error": "Unable to locate that post." }));
            } else {
                stLogger.info("sent: news post");
                res.json(newsPost);
            }
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-post server error");
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    // news sources list
    if (params.action == "news-sources" && params.lang) {
        getNews.getNewsSources(params.lang).then(newsSources => {
            stLogger.info("sent: news-sources");
            res.json(newsSources);
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-sources server error")
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    // load more
    if (params.action == "news-list-old" && params.news_id && params.sources) {
        const sources = params.sources.split(",") || [];
        getNews.getNewsList(sources, params.news_id).then(newsList => {
            stLogger.info("sent: news-list-old");
            res.json(newsList);
        }).catch(e => {
            stLogger.error("fail: news-list-old server error");
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
            console.log(e);
        });
        return;
    }

    // check for new posts
    if (params.action == "news-check" && params.sources) {
        const sources = params.sources.split(",") || [];
        getNews.getLatestId(sources).then(latestId => {
            stLogger.info("sent: news-check");
            res.json(latestId);
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-check server error");
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }


    // initial news list
    if (params.action == "news-list" && params.sources) {
        // initial news list
        const sources = params.sources.split(",") || [];
        getNews.getNewsList(sources).then(newsList => {
            stLogger.info("sent: news-list");
            res.json(newsList);
        }).catch(e => {
            console.log(e);
            stLogger.error("fail: news-check server error");
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    res.status(404).send(JSON.stringify({ "error": "Sorry can't find that!" }));
}

module.exports = {
    handle
}