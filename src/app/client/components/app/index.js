import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { RouteHandler } from 'react-router';

export default class App extends Component {

  render() {
    return <div>
      <Helmet title='Go + React = rocks!' />
      <RouteHandler {...this.props} />
    </div>;
  }

}
