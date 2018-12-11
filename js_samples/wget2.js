var http = require('http')
var options = {
    url: 'http://www.baidu.com',
    method: 'get',
};

try {
	resp = http.request(options)
	console.log(resp.statusCode)
	for (var key in resp.headers) {
		console.log(key, resp.headers[key])
	}
	console.log(resp.data)
} catch (ex) {
	console.log(ex.message)
}
