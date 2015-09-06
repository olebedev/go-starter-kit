import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { usage, todo } from './styles';
import { example, p, link } from '../homepage/styles';

export default class Usage extends Component {

  static loadProps({/* flux, state */}) {
    // Load all needed data and do what you want
    // flux.getActions('app').doSomething();
  }

  render() {
    return <div className={usage}>
      <Helmet title='Usage' />
      <h2 className={example}>Usage:</h2>
      <p className={p}>
        <span className={todo}>// TODO</span>
      </p>
      <br />
      <Link to='home' className={link}>go home</Link>
    </div>;
  }

}
