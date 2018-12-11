var url = require('url')
var r = url.parseQuery('a=b&c=d')
console.log(r)

r = url.parse('http://user:passwd@www.sample.com/path?a=b&c=d1&c=d2')
console.log(r)
