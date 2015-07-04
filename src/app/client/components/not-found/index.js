import React, { Component } from 'react';
import { notFound } from './styles';

export default class NotFound extends Component {

  static loadProps({flux, state}) {
    // Load all needed data and set the document title
    flux.getActions('app').setTitle('404 Page Not Found')
  }

  render() {
    return <div className={notFound}>
      404 Page Not Found
    </div>;
  }

}
