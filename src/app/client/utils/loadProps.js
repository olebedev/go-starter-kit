import when from 'when';

/**
 * Returns a promise that resolves after any promises returned by the routes
 * resolve. The practical uptake is that you can wait for your data to be
 * fetched before continuing. Based off react-router's async-data example
 * https://github.com/rackt/react-router/blob/master/examples/async-data/app.js
 * @param {array} components - Matched routes
 * @param {string} method - Name of static method to call
 * @param {...any} ...args - Arguments to pass to the static method
 * @return Promise
 */

export default function(components, method, ...args) {
  return when.all(components
    .map(component => component[method])
    .filter(m => typeof m === 'function')
    .map(m => m(...args))
  );
}
