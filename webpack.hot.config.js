var webpack = require('webpack');
var config = require('./webpack.config');

config.entry = {
  bundle: [
    'webpack-dev-server/client?http://localhost:5001',
    'webpack/hot/only-dev-server',
    config.entry.bundle
  ]
};

config.plugins.push(
  new webpack.HotModuleReplacementPlugin()
);

config.devServer = {
  hot: true,
  port: 5001,
  progress: true,
  publicPath: config.output.publicPath,
  stats: { colors: true },
  historyApiFallback: false,
  proxy: {'*': 'http://localhost:5000/'}
}

config.devtool = 'cheap-module-eval-source-map'

module.exports = config;
