import React, { Component } from 'react';
import { RouteHandler } from 'react-router';
import { example } from './styles';
import Homepage from '../homepage';

export default class App extends Component {

  render() {
    return <html lang='en'>
      <head>
        <meta charSet='UTF-8' />
        <link rel='stylesheet' href='/static/build/bundle.css' />
        <title></title>
      </head>
      <body>
        <RouteHandler {...this.props} />
        <script src='/static/build/bundle.js'></script>
      </body>
    </html>;
  }

}
