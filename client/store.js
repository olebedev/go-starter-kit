import { applyMiddleware, createStore as reduxCreateStore } from 'redux';
import reducers from './reducers';

const middlewares = [];

// Add state logger
if (process.env.NODE_ENV !== 'production') {
  middlewares.push(require('redux-logger')());
}

export function createStore(state) {
  return reduxCreateStore(
    reducers,
    state,
    applyMiddleware.apply(null, middlewares)
  );
}

export let store = null;
export function getStore() { return store }
export function setAsCurrentStore(s) {
  store = s;
  if (process.env.NODE_ENV !== 'production'
    && typeof window !== 'undefined') {
    window.store = store;
  }
}
