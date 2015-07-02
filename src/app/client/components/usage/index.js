import React, { Component } from 'react';
import { Link } from 'react-router';
import { example } from './styles';

export default class Usage extends Component {

  static loadProps({flux, state}) {
    // Load all needed data and set the document title
    flux.getActions('app').setTitle('Usage')
  }

  render() {
    return <div>
      Usage
      <ul>
        <li>one</li>
        <li>two</li>
        <li>three</li>
      </ul>
      <Link to='home'>home</Link>
    </div>;
  }

}
