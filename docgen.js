const fs = require('fs')
const path = require('path')
const rgx = /\$\{\d+:(.*?)\}/g

const snippets = JSON.parse(fs.readFileSync(path.join(__dirname, './snippets/snippets.json')).toString())

let out = '# Commands\n\n'

for (let snippetName in snippets) { // ToC
  if (Array.isArray(snippets[snippetName].prefix)) {
    out += `- [${snippets[snippetName].prefix[0]}](#${snippets[snippetName].prefix[0].split(' ').join('-')})\n\n`
  } else {
    out += `- [${snippets[snippetName].prefix}](#${snippets[snippetName].prefix.split(' ').join('-')})\n\n`
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

fs.writeFileSync(path.join(__dirname, 'COMMANDS.md'), out)
