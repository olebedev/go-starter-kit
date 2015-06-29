import React from 'react';
import { Route, Redirect, DefaultRoute, NotFoundRoute } from 'react-router';
import App from '#app/components/app';
import Homepage from '#app/components/homepage';
import Usage from '#app/components/usage';
import NotFound from '#app/components/not-found';


export default (
  <Route name="app" path="/" handler={App}>
    <DefaultRoute name='home' handler={Homepage} />
    <Route name="usage" path="/usage" handler={Usage} />
    {/* Server redirect in action */}
    <Redirect from="/docs" to="usage" />
    <NotFoundRoute handler={NotFound} />
  </Route>
);
