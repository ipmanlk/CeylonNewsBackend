const getNews = require("./getNews");

const handle = (req, res) => {
    const params = req.query;

    // news post
    if (params.post_id) {
        getNews.getNewsPost(params.post_id).then(newsPost => {
            res.json(newsPost);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // if sources aren't defined
    if (!params.sources || params.sources.trim() == "") {
        res.json("Sorry can't find that!");
        return;
    }

    // news sources list
    if (params.sources == "list" && params.lang) {
        getNews.getNewsSources(params.lang).then(newsSources => {
            res.json(newsSources);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // load more
    if (params.last_news_id) {
        const sources = params.sources.split(",") || [];
        getNews.getNewsList(params.lang, sources, params.last_news_id).then(newsList => {
            res.json(newsList);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // check for new posts
    if (params.lang && params.check_news && params.sources) {
        const sources = params.sources.split(",") || [];
        getNews.getLatestId(params.lang, sources).then(latestId => {
            res.json(latestId);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // initial news list
    const sources = params.sources.split(",") || [];
    getNews.getNewsList(params.lang, sources).then(newsList => {
        res.json(newsList);
    }).catch(e => {
        console.log(e);
    });
}

module.exports = {
    handle
}