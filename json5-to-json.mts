import json5 from "json5";
import fs from "node:fs";
const f = "search-podcast-episode_manual-query.json";
const v = fs.readFileSync(f, "utf8");
const vo = json5.parse(v);
const v2 = JSON.stringify(vo, null, 2);
fs.writeFileSync(f, v2);
