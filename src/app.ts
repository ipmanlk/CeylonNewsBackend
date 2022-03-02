import { startAPI } from "./api/server";
import { startScraper } from "./scraper";

await startScraper();
await startAPI();
