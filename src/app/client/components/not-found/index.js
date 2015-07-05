import React, { Component } from 'react';
import { Link } from 'react-router';
import { notFound } from './styles';
import { link } from '../homepage/styles';

export default class NotFound extends Component {

  static loadProps({flux}) {
    // Set the document title
    flux.getActions('app').setTitle('404 Page Not Found');
  }

  render() {
    return <div>
      <h2 className={notFound}>
      404 Page Not Found</h2>
      <Link to='home' className={link}>go home</Link>
    </div>;
  }

}
