import pc from "@prisma/client";

let client: any;

export function getPrismaClient(): pc.PrismaClient {
	if (!client) {
		client = new pc.PrismaClient();
	}

	return client;
}
