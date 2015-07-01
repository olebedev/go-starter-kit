import { Flux } from 'flummox';
import AppActions from './actions/app';
import AppStore from './stores/app';

export default class AppFlux extends Flux {

  constructor() {
    super();
    this.createActions('app', AppActions);
    this.createStore('app', AppStore, this);
  }

}
