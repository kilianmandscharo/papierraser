const NumericGoTypes = [
  "int",
  "int8",
  "int16",
  "int32",
  "int64",
  "uint",
  "uint8",
  "uint16",
  "uint32",
  "uint64",
  "float32",
  "float64",
];

export class StructTranspiler {
  private index: number;

  constructor(private data: string) {
    this.index = 0;
  }

  private isWhiteSpace(char: string): boolean {
    return char === "\n" || char === "\t" || char === " ";
  }

  private chopWhitespace() {
    while (
      this.isWhiteSpace(this.data[this.index]) &&
      this.index < this.data.length
    ) {
      this.index++;
    }
  }

  private chopNextWord(): string {
    this.chopWhitespace();

    const start = this.index;
    while (
      !this.isWhiteSpace(this.data[this.index]) &&
      this.index < this.data.length
    ) {
      this.index++;
    }

    return this.data.slice(start, this.index);
  }

  private parseStructSignature(): string {
    this.chopNextWord();
    const name = this.chopNextWord();
    this.chopNextWord();
    this.chopNextWord();
    return name;
  }

  private peek(): string {
    const start = this.index;
    this.chopWhitespace();
    const char = this.data[this.index];
    this.index = start;
    return char;
  }

  private parseTag(tag: string): string | null {
    tag = tag.slice(1, tag.length - 1);
    const [key, value] = tag.split(":");
    if (key !== "json") {
      return null;
    }
    return value.slice(1, value.length - 1);
  }

  private parseStructBody(): GoStructField[] {
    const fields: GoStructField[] = [];

    while (true) {
      if (this.peek() === "}") {
        break;
      }

      const key = this.chopNextWord();
      const value = this.chopNextWord();

      const field: GoStructField = {
        key,
        value,
      };

      if (this.peek() === "`") {
        const alias = this.parseTag(this.chopNextWord());
        if (alias) {
          field.alias = alias;
        }
      }

      fields.push(field);
    }

    return fields;
  }

  private parseStruct(): GoStruct {
    const name = this.parseStructSignature();
    const fields = this.parseStructBody();

    return {
      name,
      fields,
    };
  }

  private parseFieldValue(value: string): string | null {
    value = value.replace("*", "");

    if (value.startsWith("[]")) {
      return value.slice(2) + "[]";
    }

    if (NumericGoTypes.includes(value)) {
      return "number";
    }

    if (value === "bool") {
      return "boolean";
    }

    // TODO: implement maps

    return value;
  }

  private buildInterfaceString(goStruct: GoStruct): string {
    let interfaceString = `export interface ${goStruct.name} {` + "\n";

    goStruct.fields.forEach((field) => {
      interfaceString +=
        `  ${field.alias ?? field.key}: ${this.parseFieldValue(field.value)};` +
        "\n";
    });

    interfaceString += "}";

    return interfaceString;
  }

  generateInterface(): string {
    const struct = this.parseStruct();
    const interfaceString = this.buildInterfaceString(struct);
    return interfaceString;
  }
}

interface GoStructField {
  key: string;
  value: string;
  alias?: string;
}

interface GoStruct {
  name: string;
  fields: GoStructField[];
}

const structString = `type Race struct { 
	Track     Track     \`json:"track"\` 
	Players   []*Player \`json:"players"\` 
	Turn      int       \`json:"turn"\` 
	Winner    int       \`json:"winner"\` 
	GamePhase enum.GamePhase 
}`;

const structTranspiler = new StructTranspiler(structString);
console.log(structTranspiler.generateInterface());
