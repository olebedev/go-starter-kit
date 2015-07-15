import React, { Component } from 'react';
import { Link } from 'react-router';
import { example, p, link, button } from './styles';

export default class Homepage extends Component {

  static loadProps({flux}) {
    // Load all needed data and set the document title
    flux.getActions('app').setTitle('Homepage');
  }

  render() {
    const actions = this.props.flux.getActions('app');
    return <div>
      <h1 className={example}>
        Hot Reloadable <br />
        Golang + React + Flummox + Css-Module Isomorphic Starter Kit</h1>
      <br />
      <p className={p}>
        Please take a look at <Link className={link} to='/docs'>usage</Link> page.
      </p>
      {/*
      <p className={p}>
        <button onClick={actions.incrFontSize} className={button}>+</button>
        &nbsp;
        <button onClick={actions.decrFontSize} className={button}>-</button>
      </p>
      */}
    </div>;
  }

}
