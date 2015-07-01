import { Store } from 'flummox';

export default class AppStore extends Store {

  constructor(flux) {
    super();

    const appActionIds = flux.getActionIds('app');
    this.register(appActionIds.refreshStyles, this.handleRefreshStyles);

    this.state = {
      count: 0
    };
  }

  handleRefreshStyles() {
    this.setState({
      count: this.state.count + 1
    });
  }

}
