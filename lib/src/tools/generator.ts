import { resolve } from "path";
import { open, FileHandle, writeFile } from "fs/promises";

interface EnumObject {
  name: string;
  fields: string[];
}

async function parseEnums(file: FileHandle): Promise<EnumObject[]> {
  const enumObjects: EnumObject[] = [];
  let currentEnumObject: EnumObject | null = null;

  for await (let line of file.readLines()) {
    line = line.trim();

    if (line.length == 0 || line.startsWith("const (")) {
      continue;
    }

    if (line.startsWith("type ")) {
      const name = line.split(" ")[1];
      currentEnumObject = {
        name,
        fields: [],
      };
      continue;
    }

    if (line.endsWith(")") && currentEnumObject) {
      enumObjects.push(currentEnumObject);
      currentEnumObject = null;
      continue;
    }

    if (currentEnumObject) {
      const field = line.split(" ")[0];
      currentEnumObject.fields.push(field);
    }
  }

  return enumObjects;
}

function createEnumString(enumObjects: EnumObject[]): string {
  let fileString = "";

  enumObjects.forEach((enumObject, i) => {
    let enumString = `export enum ${enumObject.name} {` + "\n";
    enumObject.fields.forEach((field) => {
      enumString += `  ${field} = "${field}",` + "\n";
    });
    enumString += "}";

    fileString += enumString;

    if (i < enumObjects.length - 1) {
      fileString += "\n\n";
    }
  });

  return fileString;
}

async function generateEnums() {
  const inputfilePath = resolve(__dirname, "../../app/enum/enum.go");
  const outputFilePath = resolve(__dirname, "../src/scripts/enum.ts");

  const inputFile = await open(inputfilePath);
  const enumObjects = await parseEnums(inputFile);
  await inputFile.close();

  const output = createEnumString(enumObjects);
  await writeFile(outputFilePath, output);
}

void (async () => {
  try {
    await generateEnums();
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
})();
