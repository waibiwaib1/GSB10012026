'use strict';

var key = 'YOURAPIKEY'
  , SparkPost = require('sparkpost')
  , client = new SparkPost(key)
  , relayWebhook = {
    id: 'TEST_RELAY_WEBHOOK_UUID'
    , name: 'Updated test relay webhook'
    , target: 'http://client.test.com/updated-test-relay-webhook'
  };

client.relayWebhooks.update(relayWebhook, function(err, res) {
  if (err) {
    console.log(err);
  } else {
    console.log(res.body);
    console.log('Congrats you can use our SDK!');
  }
});
