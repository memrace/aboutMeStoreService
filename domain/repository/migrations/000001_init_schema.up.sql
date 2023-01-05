CREATE TABLE "dialogs" (
    "id" BIGINT PRIMARY KEY,
    "userName" TEXT UNIQUE,
    "firstName" TEXT,
    "lastName" TEXT,
    "chatId" INTEGER UNIQUE,
    "reply" TEXT,
    "replied" INT2
);