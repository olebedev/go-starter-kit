import React from 'react';
import { renderToString } from 'react-dom/server';
import { match, RoutingContext } from 'react-router';
import FluxComponent from 'flummox/component';
import Flux from '../flux';
import Helmet from 'react-helmet';
import routes from './routes';
import html from './html';

/**
 * Handle HTTP request at Golang server
 *
 * @param   {Object}   options  request options
 * @param   {Function} cbk      response callback
 */
export default function (options, cbk) {

  // server side fetch polyfill in action
  fetch('/api/v1/conf').then((r) => {
    return r.json();
  }).then((conf) => {

    let result = {
      error: null,
      body: null,
      redirect: null
    };


    const flux = new Flux();
    flux.getStore('app').setAppConfig(conf);

    match({ routes, location: options.url }, (error, redirectLocation, renderProps) => {
      if (error) {
        result.error = error;

      } else if (redirectLocation) {
        result.redirect = redirectLocation.pathname + redirectLocation.search;

      } else {
        const app = renderToString(
          <FluxComponent flux={flux}>
            <RoutingContext {...renderProps} />
          </FluxComponent>
        );
        const head = Helmet.rewind();
        result.body = html({app, head});
      }

      cbk(result);

    });
  });
}
