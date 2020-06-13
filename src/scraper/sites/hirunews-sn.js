const fetch = require("node-fetch");
const cheerio = require("cheerio");

const scrape = async (siteData) => {
    console.log("hirunews-sn scraper is running");
        
    const response = await fetch(siteData.site);
    const html = await response.text();
    const $ = cheerio.load(html);

    // extract post links and titles
    const posts = [];
    $(".all-section-tittle a").each((i, element) => {
        if (i == 0 || i % 2 == 0) return;
        const title = $(element).text();
        const postUrl = $(element).attr("href");

        posts.push({
            source_id: siteData.id,
            title: title,
            link: postUrl,
        });
    });

    // extract post content
    for (let p of posts) {
        const postResponse = await fetch(encodeURI(p.link));
        const postHtml = await postResponse.text();
        const $post = cheerio.load(postHtml);
        $post("#idc-container-parent").remove();
        $post("script").remove();
        $post("hr").remove();
        const postContent = $post(".main-article-section").html();
        p["post"] = postContent;
        p["main_img"] = getMainImg(postContent);
    }

    return posts;
}

const getMainImg = (post) => {
    const regEx = /<img.+src\=(?:\"|\')(.+?)(?:\"|\')(?:.+?)\>/;
    try {
        let imgs = (regEx.exec(`${post}`));
        let img = imgs[3] || imgs[2] || imgs[1];
        img = img.replace(/^http:\/\//i, 'https://');
        return img;
    } catch (error) {
        return "null";
    }
}

module.exports = {
    scrape
}