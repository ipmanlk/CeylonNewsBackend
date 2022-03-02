import express from "express";
import { getPost } from "./v1.0/handlers/get-post.handler";
import { getPosts } from "./v1.0/handlers/get-posts.handler";
import { getSources } from "./v1.0/handlers/get-sources.handler";
import cors from "cors";

export async function startAPI() {
	const app = express();
	const port = 5000;

	app.use(express.json());
	app.use(cors());

	app.get("/api/v1.0/news", getPosts);
	app.get("/api/v1.0/news/:id", getPost);
	app.get("/api/v1.0/sources", getSources);

	app.listen(port, () => {
		console.log(`API is running on port ${port}`);
	});
}
