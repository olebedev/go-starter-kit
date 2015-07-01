import React, { Component } from 'react';
import { Link } from 'react-router';
import { example } from './styles';

export default class Usage extends Component {

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
