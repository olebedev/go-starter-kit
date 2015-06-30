import React from 'react';
import Router from 'react-router';
import RenderToString from './RenderToString';
import loadProps from '#app/utils/loadProps';

import routes from './routes';

export function run() {
  Router.run(routes, Router.HistoryLocation, (Handler, state) => {
    const routeHandlerInfo = { state };
    loadProps(state.routes, 'loadProps', routeHandlerInfo).then(()=> {
      React.render(<Handler />, document);
    });
  });
};

export const renderToString = RenderToString;
