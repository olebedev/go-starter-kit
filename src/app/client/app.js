import React, { Component } from 'react';
import { example } from './styles';

export default class App extends Component {

  render() {
    return <html lang='en'>
      <head>
        <meta charSet='UTF-8' />
        <link rel='stylesheet' href='/static/build/bundle.css' />
        <title></title>
      </head>
      <body>
       <h1 className={example}>Go Starter Kit</h1>
       <script src='/static/build/bundle.js'></script>
      </body>
    </html>;
  }

}
