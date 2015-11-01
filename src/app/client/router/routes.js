import React from 'react';
import { Route, IndexRoute, Redirect } from 'react-router';
import App from '#app/components/app';
import Homepage from '#app/components/homepage';
import Usage from '#app/components/usage';
import NotFound from '#app/components/not-found';

/**
 * Returns configured routes for different
 * environments. `W` - wrapper that helps skip
 * data fetching with onEnter hook at first time.
 * @param {Object} - any data for static loaders(flux) and marker(first)
 * @returns {Object} - configured routes
 */
export default ({flux, first}) => {

  // Closure to skip firts request
  function w(loader) {
    return (nextState, replaceState, callback) => {
      if (first.time) {
        first.time = false;
        return callback();
      }
      return loader ? loader({flux, nextState, replaceState, callback}) : callback();
    };
  }

  return <Route path="/" component={App}>
    <IndexRoute component={Homepage} onEnter={w(Homepage.onEnter)}/>
    <Route path="/usage" component={Usage} onEnter={w(Usage.onEnter)}/>
    {/* Server redirect in action */}
    <Redirect from="/docs" to="/usage" />
    <Route path="*" component={NotFound} onEnter={w(NotFound.onEnter)}/>
  </Route>;
};
