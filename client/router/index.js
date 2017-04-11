import React from 'react';
import { render } from 'react-dom';
import { Router, browserHistory } from 'react-router';
import { Provider } from 'react-redux';
import toString from './toString';
import { Promise } from 'when';
import createRoutes from './routes';
import { createStore, setAsCurrentStore } from '../store';


export function run() {
    // init promise polyfill
    window.Promise = window.Promise || Promise;
    // init fetch polyfill
    window.self = window;
    require('whatwg-fetch');

    const store = createStore(window.devToolsExtension && window.devToolsExtension());
    setAsCurrentStore(store);

    // Change first.time to 'true' to make onEnter functions not to fire on the initial request

    render(
        <Provider store={store} >
            <Router history={browserHistory}>{createRoutes({store, first: { time: false }})}</Router>
        </Provider>,
        document.getElementById('app')
    );

}

// Export it to render on the Golang sever, keep the name sync with -
// https://github.com/olebedev/go-starter-kit/blob/master/src/app/server/react.go#L65
export const renderToString = toString;

require('../css');

// Style live reloading
if (module.hot) {
    let c = 0;
    module.hot.accept('../css', () => {
        require('../css');
        const a = document.createElement('a');
        const link = document.querySelector('link[rel="stylesheet"]');
        a.href = link.href;
        a.search = '?' + c++;
        link.href = a.href;
    });
}
