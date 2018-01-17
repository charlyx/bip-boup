const Https = require('https');
const { URL } = require('url');

module.exports.EmojiOrNothing = (channel, emoji) => channel.guild != null ? channel.guild.emojis.find('name', emoji) || '' : ''

module.exports.HttpsGet = (url, callback) => {
    // Need an URL object
    if (typeof(url) == typeof(''))
        url = new URL(url);

    return Https.get({
        hostname: url.hostname,
        path: url.pathname,
        headers: {
            'User-Agent': 'Bip Boup/1.0.0'
        }
    }, result => {
        result.setEncoding('utf8');

        let htmlBody = '';
        result.on('data', data => htmlBody += data);
        result.on('end', () => {
            callback(htmlBody);
        });
    });
};
