import { Store } from 'flummox';

export default class AppStore extends Store {

  constructor(flux) {
    super();

    const appActionIds = flux.getActionIds('app');
    this.register(appActionIds.refreshStyles, this.handleRefreshStyles);
    this.register(appActionIds.setTitle, this.handleSetTitle);
    this.register(appActionIds.incrFontSize, this.handleFontSize);
    this.register(appActionIds.decrFontSize, this.handleFontSize);

    this.state = {
      count: 0,
      title: 'Go + React = rocks!',
      fontSize: 150
    };
  }

  handleRefreshStyles() {
    this.setState({
      count: this.state.count + 1
    });
  }

  handleSetTitle(title) {
    this.setState({
      title: title
    });
  }

  handleFontSize(value) {
    this.setState({
      fontSize: this.state.fontSize + value * 10
    })
  }

}
