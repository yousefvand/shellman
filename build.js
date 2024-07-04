// NOTE: This script doesn't work on Windows.
// This script doesn't work for unknown reason. Use Go tool.

const fs = require("fs");
const path = require("path");

// pipe :: [Function] -> Function
const pipe = (...fns) => x => fns.reduce((v, f) => f(v), x);

// directoryFiles :: String -> [String]
function directoryFiles(dir, result = []) {
  fs.readdirSync(dir).forEach(item =>
    fs.statSync(path.join(dir, item)).isDirectory()
      ? directoryFiles(path.join(dir, item), result)
      : result.push(path.join(dir, item))
  );
  return result;
}

// objectToString :: Object -> String
const objectToString = o => JSON.stringify(o, null, 2);

// jsonFiles :: [String] -> [String]
const jsonFiles = filesArray =>
  filesArray.filter(f => f.trim().toLowerCase().slice(-5) === ".json");

// filesContent :: [String] -> [[String, Object]]
const filesContent = files =>
  files.map(f => {
    try {
      return {
        key: path.parse(f).dir.split("/")[1] + "." + path.parse(f).name,
        value: JSON.parse(fs.readFileSync(f, "utf-8")),
      };
    } catch (err) {
      throw Error(`Cannot process json file: ${f}\nError: ${err}`);
    }
  });

// arrayToObject :: [[String, Object]] -> Object
const arrayToObject = arrayOfObjects =>
  arrayOfObjects.reduce((acc, c) => {
    acc[c.key] = c.value;
    return acc;
  }, {});

// write :: String -> IO
const write = str =>
  fs.writeFileSync(path.join("snippets", "snippets.json"), str);

const app = pipe(
  directoryFiles,
  jsonFiles,
  filesContent,
  arrayToObject,
  objectToString,
  write
);

// Entry Point
app("./nsroot"); // Generate snippet file

// Generate Documentation file
// TODO:  Numerous minor parsing bugs to fix!
const snippets = JSON.parse(
  fs.readFileSync(path.join(__dirname, "./snippets/snippets.json"), "utf-8")
);

const nameSpaces = [
  ...new Set(Object.keys(snippets).map(snippet => snippet.split(".")[0])),
].sort();

structure = new Map();

for (const ns of nameSpaces) {
  members = Object.keys(snippets).filter(
    snippet => snippet.split(".")[0] === ns
  );
  structure.set(ns, members.sort());
}

const toc = ["# Commands"];
const content = [];

for (const [ns, members] of structure.entries()) {
  toc.push(`### ${ns}`);
  for (const member of members) {
    let entry = Object.entries(snippets).filter(entry => entry[0] === member);
    const snippet = entry[0];
    const snippetName = snippet[0];
    const prefix = snippet[1].prefix;
    const body = snippet[1].body;
    const description = snippet[1].description;

    // fs.writeFileSync(
    //   path.join(__dirname, "ZZZ.json"),
    //   JSON.stringify(entry[0], null, 2)
    // );

    const title = Array.isArray(prefix) ? prefix.join(" , ") : prefix;
    const link = title.split(" ").join("-");
    toc.push(`  - [${title}](#${link})`);
    content.push(`## ${title}`);
    content.push(`${description} [&uarr;](#Commands)`);

    // const rgx = /\$\{\d+:(.*?)\}/g;
    // let out = "```bash\n";
    // let code = "";
    // let y;

    // if (Array.isArray(body)) {
    //   let tmp = body.map(line => line.split("\n").join(""));
    //   tmp = tmp.map(line => line.split("\t").join("  "));
    //   tmp = tmp.map(line => unescape(line).split("\\$").join("$"));
    //   tmp = tmp.map(line => line.split(/\\\$/g).join("$"));
    //   tmp = tmp.map(line => line.split(/\\\}/g).join("}"));
    //   tmp = tmp.map(line => {
    //     let y;
    //     while ((y = rgx.exec(line)) !== null) {
    //       line = line.split(y[0]).join(y[1]);
    //     }
    //     return line;
    //   });
    //   code = tmp.join("\n");
    // } else {
    //   // single line body
    //   code = unescape(body).split("\\$").join("$");
    //   code = code.split(/\\\$/g).join("$");
    //   code = code.split("\n").join("");
    //   while ((y = rgx.exec(code)) !== null) {
    //     code = code.split(y[0]).join(y[1]);
    //   }
    //   code = code.split(/\\\}/g).join("}");
    // }
    // out += code.trim();
    // out += "\n```\n\n";
    // content.push(out);
  }
}

fs.writeFileSync(
  path.join(__dirname, "COMMANDS.md"),
  toc.concat(content).join("\n\n")
);
