import React, { Component } from 'react';
import { Link } from 'react-router';
import { example, link } from './styles';

export default class Homepage extends Component {

  static loadProps() {}

  render() {
    return <div>
      <h1 className={example}>Golang + React + Router + Flummox Isomorphic Starter Kit</h1>
      <p>Please take a look at <Link className={link} to='/docs'>usage</Link> page.</p>
    </div>;
  }

}
