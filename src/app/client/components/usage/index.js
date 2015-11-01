import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { IndexLink } from 'react-router';
import { usage, todo } from './styles';
import { example, p, link } from '../homepage/styles';

export default class Usage extends Component {

  /*eslint-disable */
  static onEnter({flux, nextState, replaceState, callback}) {
    setTimeout(() => {
      console.warn( 'Fake data loading... Edit it here --->');
      callback();
    }, 500);
  }
  /*eslint-enable */

  render() {
    return <div className={usage}>
      <Helmet title='Usage' />
      <h2 className={example}>Usage:</h2>
      <p className={p}>
        <span className={todo}>// TODO</span>
      </p>
      <br />
      go <IndexLink to='/' className={link}>home</IndexLink>
    </div>;
  }

}
