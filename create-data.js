const fs = require('fs')

const data = fs.readFileSync('./test-routing.csv')
  .toString()
  .split('"\r\n"')
  .map(row => row.split('","'))

if (!fs.existsSync('./data')){
  fs.mkdirSync('./data');
}

console.log(`writing ${data.length} files...`)
console.time('Total Time')
while(data.length > 0) {
  const [question, answer] = data.pop()
  const filename = `./data/${data.length}.json`
  fs.writeFileSync(filename, JSON.stringify({
    question, answer
  }));
}
console.timeEnd("Total Time")
