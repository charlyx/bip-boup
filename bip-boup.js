const Fs = require('fs');

const Discord = require('discord.js');
const bipboup = new Discord.Client();

console.log('Starting up...');

bipboup.on('ready', () => {
    console.log('I\'m connected to the Discord guild!');

    // Properly close connection on Ctrl-C
    process.on('SIGINT', () => {
        console.log('Shuting down...');
        bipboup.destroy().then(() => process.exit());
    });
});

Fs.readFile('.tosken', { encoding : 'utf-8' }, (err, data) => {
    if (err == null)
        bipboup.login(data.trimRight('\n'));
    else
        console.error(err.message);
});
