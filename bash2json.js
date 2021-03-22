// Convert Shellscript function into JSON array
const fs = require("fs");

// Constants
inputFile = "./test.sh";
outputFile = "./test.json";

const pipe = (...fns) => x => fns.reduce((v, f) => f(v), x);

pipe(
  // Read input file
  () => {
    try {
      return fs.readFileSync(inputFile, "utf8");
    } catch (err) {
      `Cannot read input file: ${inputFile}. Error: ${err}`;
    }
  },
  // lines
  string => string.split("\n"),
  lines =>
    // replacing leading spaces (tab size = 2: "\s\s") with "\t"
    lines.map(line => line.replace(/^\s+/, m => m.replace(/\s\s/g, "\t"))),
  // escaping \
  lines => lines.map(line => line.replace("\\", "\\")),
  // escaping $
  lines => lines.map(line => line.replace(/\$/g, "\\$")),
  // to JSON
  o => JSON.stringify(o, null, 2),
  // write output file
  data => {
    try {
      fs.writeFileSync(outputFile, data);
    } catch (error) {
      `Cannot write output file: ${outputFile}. Error: ${err}`;
    }
  }
)();
