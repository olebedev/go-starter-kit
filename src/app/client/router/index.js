import React from 'react';
import Router from 'react-router';
import RenderToString from './RenderToString';

import routes from './routes';

export function run() {
  Router.run(routes, Router.HistoryLocation, (Handler) => {
    React.render(<Handler />, document);
  });
};

export const renderToString = RenderToString;
