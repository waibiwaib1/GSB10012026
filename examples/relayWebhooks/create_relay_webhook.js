'use strict';

var key = 'YOURAPIKEY'
  , SparkPost = require('sparkpost')
  , client = new SparkPost(key)
  , relayWebhook = {
    name: 'Test relay webhook'
    , target: 'http://client.test.com/test-relay-webhook'
    , domain: 'inbox.example.com'
  };

client.relayWebhooks.create(relayWebhook, function(err, res) {
  if (err) {
    console.log(err);
  } else {
    console.log(res.body);
    console.log('Congrats you can use our SDK!');
  }
});
