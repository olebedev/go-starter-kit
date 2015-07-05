import React, { Component } from 'react';
import { Link } from 'react-router';
import { usage, todo } from './styles';
import { example, p, link } from '../homepage/styles';

export default class Usage extends Component {

  static loadProps({flux}) {
    // Load all needed data and set the document title
    flux.getActions('app').setTitle('Usage');
  }

  render() {
    return <div className={usage}>
      <h2 className={example}>Usage:</h2>
      <p className={p}>
        <span className={todo}>// TODO</span>
      </p>
      <br />
      <Link to='home' className={link}>go home</Link>
    </div>;
  }

}
