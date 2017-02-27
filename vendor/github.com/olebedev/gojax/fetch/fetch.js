
/**
 * fetch.js
 *
 * a request API compatible with window.fetch
 */

var Headers = require('node-fetch/lib/headers');

var __gofetch__ = global.__fetch__;
delete global.__fetch__;

module.exports = Fetch;

/**
 * Fetch class
 *
 * @param   Mixed    url   Absolute url or Request instance
 * @param   Object   opts  Fetch options
 * @return  Promise
 */
function Fetch(url, o) {

  // allow call as function
  if (!(this instanceof Fetch))
    return new Fetch(url, o);

  // allow custom promise
  if (!Fetch.Promise) {
    throw new Error('native promise missing, set Fetch.Promise to your favorite alternative');
  };

  if (!url) {
    throw new Error('url parameter missing');
  };

  var options = o || {};

  // wrap http.request into fetch
  return new Fetch.Promise(function(resolve, reject) {

    // normalize headers
    var headers = new Headers(options.headers || {});

    if (!headers.has('user-agent')) {
      headers.set('user-agent', 'golang-fetch/0.1 (+https://github.com/olebedev/gojax/fetch)');
    }

    headers.set('connection', 'close');

    if (!headers.has('accept')) {
      headers.set('accept', '*/*');
    }

    options.headers = headers.raw();

    // send a request
    __gofetch__(url, options, function(res){
      res.url = url;
      resolve(new Response(res));
    });
  });
};


/**
 * Response class
 *
 * @param   Object  opts  Response options
 * @return  Void
 */

function Response(r) {
  var k;
  for (k in r) {
    if (k === 'headers') {
      this[k] = new Headers(r[k]);
    } else {
      this[k] = r[k];
    }
  }
  this.ok = this.status >= 200 && this.status < 300;
}

/**
 * Decode response as json
 *
 * @return  Promise
 */
Response.prototype.json = function() {
  var _this = this;
  return new Fetch.Promise(function(resolve, reject) {
    resolve(JSON.parse(_this.body));
  });
}

/**
 * Decode response body as text
 *
 * @return  Promise
 */
Response.prototype.text = function() {
  var _this = this;
  return new Fetch.Promise(function(resolve, reject) {
    resolve(_this.body);
  });
}

Fetch.Promise = typeof Promise !== 'undefined' ? Promise : require('when').Promise;
