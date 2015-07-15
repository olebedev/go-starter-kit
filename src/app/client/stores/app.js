import { Store } from 'flummox';

export default class AppStore extends Store {

  constructor(flux) {
    super();

    const appActionIds = flux.getActionIds('app');
    this.register(appActionIds.refreshStyles, this.handleRefreshStyles);
    this.register(appActionIds.setTitle, this.handleSetTitle);
    this.register(appActionIds.incrFontSize, this.handleFontSize);
    this.register(appActionIds.decrFontSize, this.handleFontSize);
    this.register(appActionIds.loadConfig, this.handleLoadConfig);

    this.state = {
      count: 0,
      title: 'Go + React = rocks!',
      fontSize: 100,
      config: {},
      headers: {}
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
    });
  }

  setAppConfig(config) {
    this.setState({config: config});
  }

  setAppHeaders(headers) {
    this.setState({headers: headers});
  }

}
