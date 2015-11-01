import React from 'react';
import { Provider } from 'react-redux';
import { DevTools, DebugPanel, LogMonitor } from 'redux-devtools/lib/react';
import { renderToString } from 'react-dom/server';
import { match, RoutingContext } from 'react-router';
import Helmet from 'react-helmet';
import createRoutes from './routes';
import { createStore } from '../store';

/**
 * Handle HTTP request at Golang server
 *
 * @param   {Object}   options  request options
 * @param   {Function} cbk      response callback
 */
export default function (options, cbk) {

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

  try {
    match({ routes: createRoutes({store, first: { time: false }}), location: options.url }, (error, redirectLocation, renderProps) => {
      try {
        if (error) {
          result.error = error;

        } else if (redirectLocation) {
          result.redirect = redirectLocation.pathname + redirectLocation.search;

        } else {
          result.app = renderToString(
            <span>
              <Provider store={store}>
                <RoutingContext {...renderProps} />
              </Provider>
              {process.env.NODE_ENV !== 'production' && <DebugPanel top right bottom>
                <DevTools store={store} monitor={LogMonitor} visibleOnLoad={false}/>
              </DebugPanel>}
            </span>
          );
          const { title, meta } = Helmet.rewind();
          result.title = title;
          result.meta = meta;
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
