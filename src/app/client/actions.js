/**
 * action types
 */

export const SET_CONFIG = 'SET_CONFIG';

/**
 * action creators
 */

export function setConfig(config) {
  return { type: SET_CONFIG, config };
}
