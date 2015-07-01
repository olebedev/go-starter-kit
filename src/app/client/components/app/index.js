import React, { Component } from 'react';
import { RouteHandler } from 'react-router';
import Flux from 'flummox/component';
import { example } from './styles';
import Homepage from '../homepage';

class Html extends Component {
  render() {
    return <html lang='en'>
      <head>
        <meta charSet='UTF-8' />
        <link rel='stylesheet' href={'/static/build/bundle.css?' + this.props.count}/>
        <title>Count</title>
      </head>
      <body>
        {this.props.children}
        <script src='/static/build/bundle.js'></script>
      </body>
    </html>;
  }
}

export default class App extends Component {

  render() {
    return <Flux connectToStores={['app']}>
      <Html>
        <RouteHandler {...this.props} />
      </Html>
    </Flux>;
  }

}
