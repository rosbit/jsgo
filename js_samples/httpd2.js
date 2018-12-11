var http = require('http')

var server = http.createServer(function (request, response) {
	response.writeHead(200, {'Content-Type': 'application/json'})
	return {
		desc: 'json sample',
		ival: 1,
		iaval: [1, 2, 3],
		saval: ['this', 'is', 'a', 'test'],
		mval: {a: 'map', b: 'val', c: 'here'}
	}
}).listen(8888)

console.log('Server running at http://127.0.0.1:8888/')
