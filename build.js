const fs = require('fs')
const path = require('path')

// pipe :: [Function] -> Function
const pipe = (...fns) => x => fns.reduce((v, f) => f(v), x)

// directoryFiles :: String -> [String]
function directoryFiles(dir, result = []) {
  fs.readdirSync(dir).forEach(item =>
    fs.statSync(path.join(dir, item)).isDirectory() ?
      directoryFiles(path.join(dir, item), result)
      : result.push(path.join(dir, item)))
  return result
}

// objectToString :: Object -> String
const objectToString = o => JSON.stringify(o, null, 2)

// jsonFiles :: [String] -> [String]
const jsonFiles = filesArray => filesArray
  .filter(f => f.trim().toLowerCase().slice(-5) === '.json')

// filesContent :: [String] -> [[String, Object]]
const filesContent = files => files.map(f => {
  try {
    return {
      key: path.parse(f).dir.split('/')[1] + '.' + path.parse(f).name,
      value: JSON.parse(fs.readFileSync(f, 'utf-8'))
    }
  } catch (err) {
    throw Error(`Cannot process json file: ${f}\nError: ${err}`)
  }
})

// arrayToObject :: [[String, Object]] -> Object
const arrayToObject = arrayOfObjects => arrayOfObjects.reduce((acc, c) => {
  acc[c.key] = c.value
  return acc
}, {})

// write :: String -> IO
const write = str => fs.writeFileSync(path.join('snippets', 'snippets.json'), str)

const app = pipe(
  directoryFiles,
  jsonFiles,
  filesContent,
  arrayToObject,
  objectToString,
  write
)

// TODO: Refactoring
function docGen() {
  const rgx = /\$\{\d+:(.*?)\}/g
  const snippets = JSON.parse(fs.readFileSync(
    path.join(__dirname, './snippets/snippets.json'), 'utf-8'
  ))

  let out = '# Commands\n\n'
  let ns = ''

  for (let snippetName in snippets) { // Table of Contents
    const snippet = snippets[snippetName]
    if (snippetName.split('.')[0] !== ns) {
      ns = snippetName.split('.')[0]
      out += `- ${ns}\n\n`
    }
    if (Array.isArray(snippet.prefix)) {
      out += `  - [${snippet.prefix[0]}](#${snippet.prefix[0].split(' ').join('-')})\n\n`
    } else {
      out += `  - [${snippet.prefix}](#${snippet.prefix.split(' ').join('-')})\n\n`
    }
  }

  for (let snippetName in snippets) {
    const snippet = snippets[snippetName]
    out += `## \`${snippet.prefix}\`\n\n`
    out += `${snippet.description} [&uarr;](#Commands)\n\n`
    out += '```bash\n'
    let body = ''
    let y
    if (snippet.body.constructor === Array) {
      let tmp = snippet.body.map(line => line.split('\n').join(''))
      tmp = tmp.map(line => line.split('\t').join('  '))
      tmp = tmp.map(line => unescape(line).split('\\$').join('$'))
      tmp = tmp.map((line) => {
        let y
        while ((y = rgx.exec(line)) !== null) {
          line = line.split(y[0]).join(y[1])
        }
        return line
      })
      body = tmp.join('\n')
    } else { // single line body
      body = unescape(snippet.body)
      body = body.split('\n').join('')
      while ((y = rgx.exec(snippet.body)) !== null) {
        body = body.split(y[0]).join(y[1])
      }
    }
    out += body
    out += '\n```\n\n'
  }

  // Replacing \$ with $s
  out = out.replace(/\\\$/mg, '$')

  fs.writeFileSync(path.join(__dirname, 'COMMANDS.md'), out)
}

// Entry Point
app('./nsroot')

// Generate Documentation
docGen()
