require('normalize.css');
require('./global');

/**
 * Components.
 * Include all css files just if you need
 * to hot reload it. And make sure that you
 * use `webpack.optimize.DedupePlugin`
 */
require('#app/components/app/styles');
require('#app/components/homepage/styles');
require('#app/components/usage/styles');
require('#app/components/not-found/styles');
