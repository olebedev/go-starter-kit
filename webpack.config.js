var path = require('path');
var webpack = require('webpack');
var ExtractTextPlugin = require("extract-text-webpack-plugin");


var plugins = [
    new webpack.NoErrorsPlugin(),
    new webpack.optimize.DedupePlugin(),
    new ExtractTextPlugin('bundle.css')
];

// if (process.env.ANYBAR_WEBPACK === 'yep') {
//   var AnybarWebpackPlugin = require('anybar-webpack');
//   plugins.push(new AnybarWebpackPlugin());
// }
//
// var cssLoader = "css-loader?disableStructuralMinification!autoprefixer-loader!stylus-loader"
//
// if (process.env.NODE_ENV === 'production') {
//   plugins = plugins.concat([
//     new webpack.optimize.UglifyJsPlugin({output: {comments: false}}),
//     new webpack.DefinePlugin({
//       'process.env': {
//         NODE_ENV: JSON.stringify('production'),
//         BUNDLE: JSON.stringify(process.env.BUNDLE || 'bundle.min'),
//         ENV: JSON.stringify(process.env.ENV || 'prod')
//       }
//     })
//
//   ]);
//   // cssLoader = 'css-loader?disableStructuralMinification&minimize!autoprefixer-loader!stylus-loader';
// } else {
//   plugins = plugins.concat([
//     new webpack.DefinePlugin({
//       'process.env': {
//         NODE_ENV: JSON.stringify('development'),
//         BUNDLE: JSON.stringify(process.env.BUNDLE || 'bundle'),
//         ENV: JSON.stringify(process.env.ENV || 'dev')
//       }
//     })
//   ]);
// }

var cssLoader = ExtractTextPlugin.extract('style-loader', 'css-loader?disableStructuralMinification!autoprefixer-loader!stylus-loader');

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
        include: path.join(__dirname, 'src'),
        loader: 'babel'
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
};

module.exports = config;
