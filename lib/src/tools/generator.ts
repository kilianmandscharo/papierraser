import { resolve } from "path";
import { open } from "fs/promises";

async function generate() {
  const filePath = resolve(__dirname, "../src/scripts/message.ts");
  const file = await open(filePath);
  for await (const line of file.readLines()) {
    console.log(line);
  }
}

generate().catch((err) => {
  console.error(err);
  process.exit(1);
});
