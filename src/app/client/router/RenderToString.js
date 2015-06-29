import React from 'react';
import Router from 'react-router';
import routes from './routes';


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
    router.run(Handler => {
      result.Body = React.renderToString(<Handler />);
    });
  } catch (error){
    if (error.redirect) {
      result.redirect = error.redirect;
    } else {
      result.error = error;
    }
  }

  return cbk(result);
}
