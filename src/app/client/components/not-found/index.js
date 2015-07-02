import React, { Component } from 'react';

export default class NotFound extends Component {

  static loadProps({flux, state}) {
    // Load all needed data and set the document title
    flux.getActions('app').setTitle('Page Not Found')
  }

  render() {
    return <div>
      404 Not Found
    </div>;
  }

}
