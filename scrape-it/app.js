const scrapeIt = require("scrape-it");
const express = require("express");

const app = express();

app.get("/", (req, res) => {
  let url = req.query.url;
  scrapeIt(url, {})
    .then(({ status, body }) => {
      if (status !== 200) {
        res.status(status).send(body);
      }
      res.status(200).send(body);
    })
    .catch((err) => {
      res.status(500).send(err);
    });
});

app.listen(8050, () => console.log("Server running on port 8050"));
