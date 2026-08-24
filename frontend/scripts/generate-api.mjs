import { readFileSync, writeFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { convertObj } from "swagger2openapi";
import openapiTS, { astToString } from "openapi-typescript";

const specPath = fileURLToPath(new URL("../../backend/docs/swagger.json", import.meta.url));
const outPath = fileURLToPath(new URL("../src/types/api-generated.ts", import.meta.url));

const spec = JSON.parse(readFileSync(specPath, "utf8"));
const { openapi } = await convertObj(spec, { patch: true, warnOnly: true });
const ast = await openapiTS(openapi);
writeFileSync(outPath, astToString(ast));
console.log(`Generated ${outPath}`);
