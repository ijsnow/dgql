const path = require('path');

module.exports = ({ config, mode }) => {
  config.module.rules.push({
    test: /\.(ts|tsx)$/,
    loader: require.resolve('babel-loader'),
    options: {
      presets: [['next/babel', { flow: false, typescript: true }]],
    },
  });

  config.resolve.extensions.push('.ts', '.tsx');

  config.module.rules.push({
    test: /\.css$/,
    loaders: [
      // Loader for webpack to process CSS with PostCSS
      {
        loader: 'postcss-loader',
        options: {
          /* 
            Enable Source Maps
           */
          sourceMap: true,
          /*
            Set postcss.config.js config path && ctx 
           */
          config: {
            path: './.storybook/',
          },
        },
      },
    ],

    include: path.resolve(__dirname, '../'),
  });

  return config;
};
