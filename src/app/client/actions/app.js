import { Actions } from 'flummox';

export default class AppActions extends Actions {

  refreshStyles() {
    return 1;
  }

  setTitle(title) {
    return title;
  }

  incrFontSize() {
    return 1;
  }

  decrFontSize() {
    return -1;
  }

  loadConfig() {
    return 'loading';
  }

}
