import { Store } from 'flummox';

export default class AppStore extends Store {

  constructor() {
    super();

    this.state = {
      config: {}
    };
  }

  setAppConfig(config) {
    this.setState({config: config});
  }

}
