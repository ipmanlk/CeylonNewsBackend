import { startAPI } from "./api/server.js";
import { startScraper } from "./scraper.js";

await startScraper();
await startAPI();
