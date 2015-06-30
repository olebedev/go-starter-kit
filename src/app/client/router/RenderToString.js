import React from 'react';
import Router from 'react-router';
import routes from './routes';
import loadProps from '#app/utils/loadProps';

/**
 * Handle HTTP request at Golang server
 *
 * @param   Object   options  request options
 * @param   Function cbk      response callback
 */
export default function (options, cbk) {

  let result = {
    error: null,
    body: null,
    redirect: null
  };


  const router = Router.create({
    routes: routes,
    location: options.url,
    onError: error => {
      throw error;
    },
    onAbort: abortReason => {
      const error = new Error();

      if (abortReason.constructor.name === 'Redirect') {
        const { to, params, query } = abortReason;
        const url = router.makePath(to, params, query);
        error.redirect = url;
      }

      throw error;
    }
  });

  try {
    router.run((Handler, state) => {
      const routeHandlerInfo = { state };
      loadProps(state.routes, 'loadProps', routeHandlerInfo).then(()=> {
        result.body = React.renderToString(<Handler />);
        cbk(result);
      });
    });
  } catch (error){
    if (error.redirect) {
      result.redirect = error.redirect;
    } else {
      result.error = error;
    }
    cbk(result);
  }

}
