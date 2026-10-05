'use strict';

var api = 'relay-webhooks'
  , toApiFormat = require('./toApiFormat');

module.exports = function(client) {
  var relayWebhooks = {
    all: function(callback) {
      var options = {
        uri: api
      };

      client.get(options, callback);
    },
    describe: function(options, callback) {
      options = options || {};

      if(!options.id) {
        callback(new Error('id is required'));
        return;
      }

      var reqOpts = {
        uri: api + '/' + options.id
      };

      client.get(reqOpts, callback);
    },
    create: function(relayWebhook, callback) {
      if(typeof relayWebhook === 'function') {
        callback = relayWebhook;
        relayWebhook = null;
      }

      if(!relayWebhook) {
        callback(new Error('relayWebhook object is required'));
        return;
      }

      var options = {
        uri: api
        , json: toApiFormat(relayWebhook)
      };

      client.post(options, callback);
    },
    update: function(relayWebhook, callback) {
      if(typeof relayWebhook === 'function') {
        callback = relayWebhook;
        relayWebhook = null;
      }

      if(!relayWebhook) {
        callback(new Error('relayWebhook object is required'));
        return;
      }

      var object = toApiFormat(relayWebhook)
        , options = {
          uri: api + '/' + relayWebhook.id
          , json: object
        };

      client.put(options, callback);
    }
  };

  relayWebhooks['delete'] = function(id, callback) {
    if (typeof id === 'function') {
      callback = id;
      id = null;
    }

    if (!id) {
      callback(new Error('id is required'));
      return;
    }

    var options = {
      uri: api + '/' + id
    };

    client['delete'](options, callback);
  };

  return relayWebhooks;
};
