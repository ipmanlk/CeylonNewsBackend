const express = require('express');
const app = express();
const port = 3000;

app.use(function (req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    next();
});

app.get('/:version', (req, res) => {
    try {
        const api = require(`./api/${req.params.version}/handle.js`);
        api.handle(req, res);
    } catch (error) {
        console.log(error);
        res.status(500).json({ "error": "Internal Server Error." });
    }
});

app.use((req, res, next) => {
    res.status(404).json({ "error": "Sorry can't find that!" });
});

app.listen(port, () => {
    console.log(`Ceylon News is running on port ${port}!`);
    const scrapper = require("./scrapper/scrape");
    console.log("Preparing the database.");
    scrapper.prepareDB().then(() => {
        console.log("Scrape CronJob started!.");
        scrapper.scrapeCronJob.start();
    }).catch(e => {
        console.log(e);
    });
});

// for testing
module.exports = app; 