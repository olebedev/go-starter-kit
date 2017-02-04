import { combineReducers } from 'redux';
import { SET_CONFIG } from './actions';

function config(state = {}, action) {
  switch (action.type) {
  case SET_CONFIG:
    return action.config;
  default:
    return state;
  }
}

export default combineReducers({config});
