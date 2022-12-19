import Koa from "koa";
import Router from "@koa/router";
import { koaBody } from "koa-body";
import amqp from "amqplib";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();
const conn = await amqp.connect(process.env.RABBITMQ_URL || "amqp://localhost");
const channel = await conn.createChannel();

const app = new Koa();
const router = new Router();

router.post("/exec", koaBody(), async (ctx) => {
  const body = ctx.request.body;
  if (!(body.code && body.language)) {
    return ctx.body = {
      id: null,
      message: "Invalid request"
    }
  }

  const submission = await prisma.submission.create({
    data: {
      code: body.code,
      input: body.input ?? "",
      language: body.language,
    },
  })

  channel.sendToQueue(
    "exec_queue",
    Buffer.from(
      JSON.stringify({
        id: submission.id,
        code: body.code,
        input: body.input,
        language: body.language,
      }),
    ),
    {
      correlationId: submission.id,
      replyTo: "exec_queue",
    },
  );

  ctx.body = {
    id: submission.id,
    message: "Submission received"
  }
});

router.get("/exec/:id", async (ctx) => {
  const id = ctx.params.id;
  const submission = await prisma.submission.findUnique({
    where: {
      id,
    },
  });
  ctx.body = submission;
});

channel.consume("exec_queue", async (msg) => {
  if (!msg) return;

  const data = JSON.parse(msg.content.toString());
  await prisma.submission.update({
    where: {
      id: msg.properties.correlationId,
    },
    data: {
      status: "Executed",
      stdout: data.stdout,
      stderr: data.stderr,
      exitCode: data.exitCode,
      message: data.message,
    },
  });
});

app.use(router.routes());
app.listen(3004);
