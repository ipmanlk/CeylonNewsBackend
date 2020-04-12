const newsData = require("./data/newsData");

const handle = (req, res) => {
    const action = req.query.action;

    // search & navigate news list 
    if (action == "news-list") {
        const sources = req.query.sources.split(",") || [];
        const keyword = req.query.keyword || "";
        const skip = req.query.skip || 0;
        newsData.getNewsList(sources, keyword, skip).then(data => {
            res.json(data);
        }).catch(e => {
            console.log(e);
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    // retrieve single news post
    if (action == "news-post") {
        getNews.getNewsPost(req.query.news_id).then(newsPost => {
            if (newsPost == undefined) {
                res.status(404).send(JSON.stringify({ "error": "Unable to locate that post." }));
            } else {
                res.json(newsPost);
            }
        }).catch(e => {
            console.log(e);
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    // retrieve news sources list
    if (action == "sources-list" && req.query.lang) {
        newsData.getNewsSources(req.query.lang).then(newsSources => {
            res.json(newsSources);
        }).catch(e => {
            console.log(e);
            res.status(500).send(JSON.stringify({ "error": "Internal Server Error." }));
        });
        return;
    }

    // default response
    res.status(404).send(JSON.stringify({ "error": "Sorry can't find that!" }));

}

module.exports = {
    handle
}