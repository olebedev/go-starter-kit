var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require("extract-text-webpack-plugin");


var plugins = [
    new webpack.NoErrorsPlugin(),
    new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin('bundle.css')
];

var cssLoader = ExtractTextPlugin.extract('style-loader', 'css-loader?disableStructuralMinification!autoprefixer-loader!stylus-loader?paths=src/app/client/styles/&import=./ctx');

var config  = {
  entry: [path.join(__dirname, 'src/app/client/entry.js')],
  output: {
    path: path.join(__dirname, 'src/app/server/data/static/build'),
    publicPath: "/static/build",
    filename: 'bundle.js'
  },
  plugins: plugins,
  module: {
    loaders: [
      {test: /\.styl$/, loader: cssLoader},
      {test: /\.(png|gif)$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000'},
      {test: /\.svg$/, loader: 'url-loader?name=[name]@[hash].[ext]&limit=5000!svgo-loader?useConfig=svgo1'},
      {test: /\.(pdf|ico|jpg|eot|otf|woff|ttf|mp4|webm)$/, loader: 'file-loader?name=[name]@[hash].[ext]'},
      {test: /\.json$/, loader: 'json-loader'},
      {
        test: /\.jsx?$/,
        include: path.join(__dirname, 'src/app/client'),
        loaders: ['babel']
      }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.jsx', '.styl'],
    alias: {
      '#app': path.join(__dirname, '/src/app/client')
    }
  },
  svgo1: {
    multipass: true,
    plugins: [
      // by default enabled
      {mergePaths: false},
      {convertTransform: false},
      {convertShapeToPath: false},
      {cleanupIDs: false},
      {collapseGroups: false},
      {transformsWithOnePath: false},
      {cleanupNumericValues: false},
      {convertPathData: false},
      {moveGroupAttrsToElems: false},
      // by default disabled
      {removeTitle: true},
      {removeDesc: true}
    ]
  },
  stylus: {
    // paths: [path.join(__dirname, 'src/app/client/styles')] // ,
    // import: ['./ctx']
  }
};

module.exports = config;
