const sqlite3 = require('sqlite3').verbose();

class SQLite {
    constructor(dbPath) {
        this.db = new sqlite3.Database(dbPath, sqlite3.OPEN_READWRITE, (err) => {
            if (err) {
                console.error(err.message);
            }
        });
    }

    // this will only get the first row in result set
    getOne(sql, args = []) {
        return new Promise((resolve, reject) => {
            this.db.get(sql, args, (err, row) => {
                if (err) {
                    reject(err);
                }
                resolve(row);
            });
        });
    }

    // this will get the entire result set
    getAll(sql, args = []) {
        return new Promise((resolve, reject) => {
            this.db.all(sql, args, (err, rows) => {
                if (err) {
                    reject(err);
                }
                resolve(rows);
            });
        });
    }

    // this will run any given query in db
    run(sql, args = []) {
        return new Promise((resolve, reject) => {
            this.db.run(sql, args, function (err) {
                if (err) {
                    reject(err);
                }
                resolve(this.lastID || true);
            });
        });
    }
}
module.exports = {
    SQLite
}