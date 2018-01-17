const { HttpsGet } = require('../utils.js');
const HtmlParser = require('htmlparser2');

module.exports.command = 'nightcore';
module.exports.help = 'Recherche du nightcore sur YouTube.';

const URL = query => `https://www.youtube.com/results?sp=EgIQAQ%253D%253D&search_query=${query}`;
const RepliedURL = query => `https://www.youtube.com${query}`;
module.exports.callback = (message, words) => {
    if (words.length > 1) {
        HttpsGet(URL(words.join('+')), (htmlBody) => {
            let done = false;
            const parser = new HtmlParser.Parser({
                onopentag: (name, attribs) => {
                    if (done) return;
                    if (name === 'a' && attribs.href.indexOf('/watch') === 0) {
                        // Somewhat filter the results
                        if (attribs.title != null && attribs.title.toLowerCase().indexOf('nightcore') >= 0) {
                            message.reply(RepliedURL(attribs.href));
                            done = true;
                        }
                    }
                }
            });
            parser.write(htmlBody);
            parser.end();

            if (!done)
                message.reply('Je n\'ai rien trouvé de satisfaisant :frowning:')
        });
    } else {
        message.reply('Tu cherches quoi ?');
    }
};
