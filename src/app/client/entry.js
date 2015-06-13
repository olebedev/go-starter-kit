var React = require('expose?React!react');
var App = require('expose?App!./app');

if(typeof window !== 'undefined') {
  React.render(<App />, document);
};

// if (process.env.ENV !==  'prod')
//   require('#js/dev-tools/reloader');
