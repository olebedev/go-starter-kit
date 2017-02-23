import React from 'react';
import { Provider } from 'react-redux';
import { renderToString } from 'react-dom/server';
import { match, RouterContext } from 'react-router';
import Helmet from 'react-helmet';
import createRoutes from './routes';
import { createStore, setAsCurrentStore } from '../store';

/**
 * Handle HTTP request at Golang server
 *
 * @param   {Object}   options  request options
 * @param   {Function} cbk      response callback
 */
export default function (options, cbk) {
  cbk = global[cbk];
  let result = {
    uuid: options.uuid,
    app: null,
    title: null,
    meta: null,
    initial: null,
    error: null,
    redirect: null
  };

  const store = createStore();
  setAsCurrentStore(store);

  try {
    match({ routes: createRoutes({store, first: { time: false }}), location: options.url }, (error, redirectLocation, renderProps) => {
      try {
        if (error) {
          result.error = error;

        } else if (redirectLocation) {
          result.redirect = redirectLocation.pathname + redirectLocation.search;

        } else {
          result.app = renderToString(
            <Provider store={store}>
              <RouterContext {...renderProps} />
            </Provider>
          );
          const { title, meta } = Helmet.rewind();
          result.title = title.toString();
          result.meta = meta.toString();
          result.initial = JSON.stringify(store.getState());
        }
      } catch (e) {
        result.error = e;
      }
      return cbk(result);
    });
  } catch (e) {
    result.error = e;
    return cbk(result);
  }
}
