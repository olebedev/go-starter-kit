import React from 'react';
import Router from 'react-router';
import FluxComponent from 'flummox/component';
import Flux from '../flux';
import RenderToString from './RenderToString';
import loadProps from '#app/utils/loadProps';

import routes from './routes';

const flux = new Flux();

export function run() {
  // share flux instance
  window.flux = flux;

  fetch('/api/v1/conf').then((r) => {
    return r.json();
  }).then((conf) => {

    flux.getStore('app').setAppConfig(conf);

    Router.run(routes, Router.HistoryLocation, (Handler, state) => {
      const routeHandlerInfo = { flux, state };
      loadProps(state.routes, 'loadProps', routeHandlerInfo).then(()=> {
        React.render(
          <FluxComponent flux={flux}>
            <Handler />
          </FluxComponent>,
          document
        );
      });
    });

  });
}

export const renderToString = RenderToString;

require('../styles');

// Style live reload
if (module.hot) {
  const refreshStyles = flux.getActions('app').refreshStyles;
  module.hot.accept('../styles', () => {
    require('../styles');
    refreshStyles();
  });
}
