const path = require('path');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  mode: 'production',
  entry: {
    index: ['./web_src/js/index', './web_src/js/draw']
  },
  devtool: 'source-map',
  output: {
    path: path.resolve(__dirname, 'public/js'),
    filename: 'index.js'
  },
  optimization: {
    minimize: true,
    minimizer: [new TerserPlugin({
      sourceMap: true,
    })],
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              [
                '@babel/preset-env',
                {
                  useBuiltIns: 'entry',
                  corejs: 3,
                }
              ]
            ]
          }
        }
      }
    ]
  }
};
