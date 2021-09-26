module.exports = {
    publicPath: process.env.NODE_ENV === 'production' ? 'https://asoul.cdn.n3ko.co/' : '/',
    transpileDependencies: [
        'vuetify'
    ]
}
