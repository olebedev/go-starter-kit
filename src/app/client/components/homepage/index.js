import React, { Component } from 'react';
import { Link } from 'react-router';
import { example, link } from './styles';

export default class Homepage extends Component {

  static loadProps() {}

  constructor() {
    super();
    this.onClick = this.onClick.bind(this);
  }

  onClick() {
    this.props.flux.getActions('app').refreshStyles();
  }

  render() {
    return <div>
      <h1 className={example}>Golang + React + Router + Redux Isomorphic Starter Kit</h1>
      <p>Please take a look at <Link className={link} to='/docs'>usage</Link> page.</p>
      <button type='button' onClick={this.onClick}>+</button>
    </div>;
  }

}
