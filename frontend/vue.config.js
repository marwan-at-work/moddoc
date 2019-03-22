module.exports = {
    publicPath: process.env.NODE_ENV === 'production' ? '/public/' : '/',
    devServer: {
        proxy: 'http://localhost:3001'
    }
}