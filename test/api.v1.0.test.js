process.env.NODE_ENV = 'test';

const chai = require("chai");
const chaiHttp = require("chai-http");
const app = require("../server");

chai.use(chaiHttp);
chai.should();

// news sources
describe("News Sources", () => {
    describe("GET /v1.0/?action=news-sources&lang=en", () => {
        it("should get an array with all English news sources with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-sources&lang=en")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("array").to.not.be.empty;
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-sources&lang=123", () => {
        it("should get an empty array with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-sources&lang=123")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("array").to.be.empty;
                    done();
                });
        });
    });
});


// news list
describe("News List", () => {
    describe("GET /v1.0/?action=news-list&sources=1,2,3,4", () => {
        it("should get an non empty array with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-list&sources=1,2,3,4")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("array").to.not.be.empty;
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-list&sources=1,2,3,4", () => {
        it("should get an non empty array with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-list&sources=1,2,3,4")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("array").to.not.be.empty;
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-list&sources=4", () => {
        it("should get an empty array with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-list&sources=4")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("array").to.not.be.empty;
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-list&sources=", () => {
        it("should get an error object with 404 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-list&sources=")
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a("object");
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-list", () => {
        it("should get an error object with 404 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-list")
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a("object");
                    done();
                });
        });
    });
});

// news post
describe("News Post", () => {
    describe("GET /v1.0/?action=news-post&post_id=1", () => {
        it("should get an object with 200 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-post&post_id=1")
                .end((err, res) => {
                    res.should.have.status(200);
                    res.body.should.be.a("object");
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-post&post_id=99999", () => {
        it("should get an error object with 404 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-post&post_id=99999")
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a("object");
                    done();
                });
        });
    });

    describe("GET /v1.0/?action=news-post", () => {
        it("should get an error object with 404 status", (done) => {
            chai.request(app)
                .get("/v1.0/?action=news-post")
                .end((err, res) => {
                    res.should.have.status(404);
                    res.body.should.be.a("object");
                    done();
                });
        });
    });
});