const getNews = require("./getNews");

const handle = (req, res) => {
    const params = req.query;

    // news post
    if (params.action == "news-post" && params.post_id) {
        getNews.getNewsPost(params.post_id).then(newsPost => {
            if (newsPost == undefined) {
                res.status(404).send(JSON.stringify({ "error": "Unable to locate that post." }));
            } else {
                res.json(newsPost);
            }
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // news sources list
    if (params.action == "news-sources" && params.lang) {
        getNews.getNewsSources(params.lang).then(newsSources => {
            res.json(newsSources);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // load more
    if (params.action == "news-list-old" && params.news_id && params.sources) {
        const sources = params.sources.split(",") || [];
        getNews.getNewsList(sources, params.news_id).then(newsList => {
            res.json(newsList);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    // check for new posts
    if (params.action == "news-check" && params.sources) {
        const sources = params.sources.split(",") || [];
        getNews.getLatestId(sources).then(latestId => {
            res.json(latestId);
        }).catch(e => {
            console.log(e);
        });
        return;
    }


    // initial news list
    if (params.action == "news-list" && params.sources) {
        // initial news list
        const sources = params.sources.split(",") || [];
        getNews.getNewsList(sources).then(newsList => {
            res.json(newsList);
        }).catch(e => {
            console.log(e);
        });
        return;
    }

    res.status(404).send(JSON.stringify({ "error": "Sorry can't find that!" }));
}

module.exports = {
    handle
}