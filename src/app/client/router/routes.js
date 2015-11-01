import React from 'react';
import { Route, IndexRoute, Redirect } from 'react-router';
import App from '#app/components/app';
import Homepage from '#app/components/homepage';
import Usage from '#app/components/usage';
import NotFound from '#app/components/not-found';


export default <Route path="/" component={App}>
  <IndexRoute component={Homepage} onEnter={Homepage.onEnter}/>
  <Route path="/usage" component={Usage} />
  {/* Server redirect in action */}
  <Redirect from="/docs" to="/usage" />
  <Route path="*" component={NotFound} />
</Route>;
