-- CreateTable
CREATE TABLE "Post" (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "title" TEXT COLLATE NOCASE NOT NULL,
    "content" TEXT COLLATE NOCASE NOT NULL,
    "url" TEXT NOT NULL,
    "thumbnailUrl" TEXT,
    "language" TEXT NOT NULL,
    "sourceName" TEXT COLLATE NOCASE NOT NULL,
    "createdDate" DATETIME NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "Post_url_key" ON "Post"("url");
