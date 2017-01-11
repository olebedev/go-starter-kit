
/**
 * fetch.js
 *
 * a request API compatible with window.fetch
 */

var Headers = require('node-fetch/lib/headers');
var assign = require('lodash/object/assign');


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
      headers.set('user-agent', 'golang-fetch/0.0 (+https://github.com/olebedev/go-duktape-fetch)');
    }

    headers.set('connection', 'close');

    if (!headers.has('accept')) {
      headers.set('accept', '*/*');
    }

    options.headers = headers.raw();

    // send a request
    var res = Fetch.goFetchSync(url, options);
    res.url = url;

    resolve(new Response(res));
  });
};


/**
 * Response class
 *
 * @param   Object  opts  Response options
 * @return  Void
 */

function Response(r) {
  assign(this, r)
  this.ok = this.status >= 200 && this.status < 300;
}

/**
 * Decode response as json
 *
 * @return  Promise

 */
Response.prototype.json = function() {

  return this.text().then(function(text) {
    return JSON.parse(text);
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
