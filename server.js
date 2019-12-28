const express = require('express')
const app = express()
const port = 3000

app.get('/:version', (req, res) => {
    const api = require(`./api/${req.params.version}/handle.js`);
    api.handle(req, res);
});

app.use((req, res, next) => {
    res.status(404).send("Sorry can't find that!")
});

app.listen(port, () => console.log(`Ceylon News is running on port ${port}!`));