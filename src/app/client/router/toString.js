import React from 'react';
import { renderToString } from 'react-dom/server';
import { match, RoutingContext } from 'react-router';
import FluxComponent from 'flummox/component';
import Flux from '../flux';
import Helmet from 'react-helmet';
import routes from './routes';

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
      uuid: options.uuid,
      app: null,
      title: null,
      meta: null,
      initial: null,
      error: null,
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
        // load data into the flux instance
        try {
          result.initial = flux.serialize();
        } catch (e) {
          result.error = 'serialization error: ' + e
        }

        result.app = renderToString(
          <FluxComponent flux={flux}>
            <RoutingContext {...renderProps} />
          </FluxComponent>
        );
        const { title, meta } = Helmet.rewind();
        result.title = title;
        result.meta = meta;
      }

      cbk(result);

    });
  });
}
