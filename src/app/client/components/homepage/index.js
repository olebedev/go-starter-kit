import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import { example, p, link } from './styles';

export default class Homepage extends Component {

  static onEnter(nextState, replaceState, callback) {
    callback()
    // setTimeout(() => {
    //   callback()
    // }, 1)
  }

  render() {
    return <div>
      <Helmet
        title='Home page'
        meta={[
          {
            property: 'og:title',
            content: 'Golang Isomorphic React/Hot Reloadable/Flummox/Css-Modules Starter Kit'
          }
        ]} />
      <h1 className={example}>
        Hot Reloadable <br />
        Golang + React + Flummox + Css-Modules
        <br />Isomorphic Starter Kit</h1>
      <br />
      <p className={p}>
        Please take a look at <Link className={link} to='/docs'>usage</Link> page.
      </p>
    </div>;
  }

}
