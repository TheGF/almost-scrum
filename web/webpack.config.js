const HtmlWebPackPlugin = require("html-webpack-plugin");
const path = require('path');

module.exports = {
    devServer: {
        contentBase: './dist',
        hot: true,
        proxy: {
            '/api': 'http://127.0.0.1:4444',
        },  
        watchOptions: {
            ignored: [
              path.resolve(__dirname, 'dist'),
              path.resolve(__dirname, 'node_modules')
            ]
        }       
    },
    devtool: 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                resolve: { extensions: [".js", ".jsx"] },
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader"
                }
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader"
                    }
                ]
            },
            // {
            //     test: /\.css$/i,
            //     use: ['style-loader', 'css-loader'],
            // },
            {
                test: /\.s[ac]ss$/i,
                use: ['style-loader', 'css-loader', 'sass-loader'],
              },
        ]
    },
    plugins: [
        new HtmlWebPackPlugin({
            template: "./src/index.html",
            filename: "./index.html"
        })
    ],
};
