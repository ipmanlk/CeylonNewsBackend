const fetch = require("node-fetch");
const cheerio = require("cheerio");

const scrape = async (siteData) => {
  console.log("adaderana-sn scraper: started");

  const response = await fetch(siteData.site);
  const html = await response.text();
  const $ = cheerio.load(html);

  // extract post links and titles
  const posts = [];
  $(".story-text").each((i, storyText) => {

    // only load 5 posts at once
    if (i > 4) return false;

    const a = $(storyText).find("h2 a");
    const title = $(a).text();
    const postUrl = "http://sinhala.adaderana.lk/" + $(a).attr("href");

    posts.push({
      source_id: siteData.id,
      title: title,
      link: postUrl,
    });

  });

  // extract post content
  for (let p of posts) {
    const postResponse = await fetch(encodeURI(`https://s1.navinda.xyz/splash/render.html?url=${p.link}&timeout=13&wait=8`));
    const postHtml = await postResponse.text();
    const $post = cheerio.load(postHtml);

    $post(".uiGrid _51mz").remove();
    $post("script").remove();
    $post("hr").remove();

    const postContent = $post(".news-content").html();
    p["post"] = postContent;
    p["main_img"] = getMainImg($post.html());
  }

  console.log("adaderana-sn scraper: stopped");
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