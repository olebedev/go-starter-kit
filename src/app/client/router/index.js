import React from 'react';
import Router from 'react-router';
import FluxComponent from 'flummox/component';
import Flux from '../flux';
import RenderToString from './RenderToString';
import loadProps from '#app/utils/loadProps';

import routes from './routes';

const flux = new Flux();

export function run() {
  Router.run(routes, Router.HistoryLocation, (Handler, state) => {
    const routeHandlerInfo = { state };
    loadProps(state.routes, 'loadProps', routeHandlerInfo).then(()=> {
      React.render(
        <FluxComponent flux={flux}>
          <Handler />
        </FluxComponent>,
        document
      );
    });
  });
};

export const renderToString = RenderToString;
