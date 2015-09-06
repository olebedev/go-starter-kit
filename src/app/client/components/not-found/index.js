import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { notFound } from './styles';
import { link } from '../homepage/styles';

export default class NotFound extends Component {

  static loadProps({/* flux, state */}) {
    // Load all needed data and do what you want
    // flux.getActions('app').doSomething();
  }

  render() {
    return <div>
      <Helmet title='404 Page Not Found' />
      <h2 className={notFound}>
      404 Page Not Found</h2>
      <Link to='home' className={link}>go home</Link>
    </div>;
  }

}
