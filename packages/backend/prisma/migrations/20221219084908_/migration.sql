-- CreateTable
CREATE TABLE "Submission" (
    "id" TEXT NOT NULL,
    "language" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'In Queue',
    "message" TEXT NOT NULL DEFAULT '',
    "code" TEXT NOT NULL,
    "input" TEXT NOT NULL,
    "stdout" TEXT NOT NULL DEFAULT '',
    "stderr" TEXT NOT NULL DEFAULT '',
    "exitCode" INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT "Submission_pkey" PRIMARY KEY ("id")
);
