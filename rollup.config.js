import babel from '@rollup/plugin-babel';
import replace from '@rollup/plugin-replace';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import postcss from 'rollup-plugin-postcss';
import externalGlobals from 'rollup-plugin-external-globals';
import svgr from '@svgr/rollup';

export default {
  input: 'src/index.js',
  output: {
    dir: 'dist',
    entryFileNames: 'extension.[hash].js',
    format: 'esm'
  },
  plugins: [
    resolve({
      customResolveOptions: {
        moduleDirectory: 'node_modules'
      }
    }),
    babel({
      plugins: [
        '@babel/plugin-proposal-object-rest-spread',
        '@babel/plugin-proposal-optional-chaining',
        '@babel/plugin-syntax-dynamic-import',
        '@babel/plugin-proposal-class-properties',
        'transform-react-remove-prop-types',
      ],
      exclude: 'node_modules/**'
    }),
    commonjs({
      include: 'node_modules/**',
      namedExports: {
        'node_modules/react-is/index.js': ['ForwardRef', 'isValidElementType', 'isElement'],
        'node_modules/react/index.js': ['createElement']
      }
    }),
    postcss(),
    externalGlobals({
      'carbon-components-react': 'CarbonComponentsReact',
      'react': 'React',
      'react-dom': 'ReactDOM',
      'react-redux': 'ReactRedux',
      'react-router-dom': 'ReactRouterDOM',
      'stream': 'Stream'
    }),
    replace({
      'process.env.NODE_ENV': JSON.stringify( 'production' )
    }),
    svgr(),
  ]
};