const express = require("express");
const cors = require("cors");
const app = express();
const port = 3000;

app.use(cors());

app.get("/:version", (req, res) => {
    try {
        const api = require(`./api/${req.params.version}/handle.js`);
        api.handle(req, res);
    } catch (error) {
        console.log(error);
        res.status(404).json({ "error": "Sorry!. I can't find that." });
    }
});

app.use((req, res, next) => {
    res.status(404).json({ "error": "Sorry!. I can't find that!" });
});

app.listen(port, () => {
    console.log(`Ceylon News is running on port ${port}!`);
    const scraper = require("./scraper/scrape");
    console.log("Preparing the database.");
    scraper.prepareDB().then(() => {
        console.log("Scrape CronJob started!.");
        scraper.scrapeCronJob.start();
    }).catch(e => {
        console.log(e);
    });
});

// for testing
module.exports = app; 