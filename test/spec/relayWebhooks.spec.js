var chai = require('chai')
  , expect = chai.expect
  , sinon = require('sinon')
  , sinonChai = require('sinon-chai');

chai.use(sinonChai);

describe('Relay Webhooks Library', function() {
  var client, relayWebhooks;

  beforeEach(function() {
    client = {
      get: sinon.stub().yields(),
      post: sinon.stub().yields(),
      put: sinon.stub().yields(),
      'delete': sinon.stub().yields()
    };

    relayWebhooks = require('../../lib/relayWebhooks')(client);
  });

  describe('all Method', function() {
    it('should call client get method with the appropriate uri', function(done) {
      relayWebhooks.all(function(err, data) {
        expect(client.get.firstCall.args[0].uri).to.equal('relay-webhooks');
        done();
      });
    });
  });

  describe('describe Method', function() {
    it('should call client get method with the appropriate uri', function(done) {
      var options = {
        id: 'test'
      };

      relayWebhooks.describe(options, function(err, data) {
        expect(client.get.firstCall.args[0].uri).to.equal('relay-webhooks/test');
        done();
      });
    });

    it('should throw an error if id is missing', function(done) {
      relayWebhooks.describe(null, function(err) {
        expect(err.message).to.equal('id is required');
        expect(client.get).not.to.have.been.called;
        done();
      });
    });
  });

  describe('create Method', function() {
    it('should call client post method with the appropriate uri', function(done) {
      relayWebhooks.create({}, function(err, data) {
        expect(client.post.firstCall.args[0].uri).to.equal('relay-webhooks');
        done();
      });
    });

    it('should throw an error if relayWebhook is null', function(done) {
      relayWebhooks.create(null, function(err) {
        expect(err.message).to.equal('relayWebhook object is required');
        expect(client.post).not.to.have.been.called;
        done();
      });
    });

    it('should throw an error if relayWebhook is missing', function(done) {
      relayWebhooks.create(function(err) {
        expect(err.message).to.equal('relayWebhook object is required');
        expect(client.post).not.to.have.been.called;
        done();
      });
    });
  });

  describe('update Method', function() {
    it('should call client put method with the appropriate uri', function(done) {
      var relayWebhook = {
        id: 'test'
      };

      relayWebhooks.update(relayWebhook, function(err, data) {
        expect(client.put.firstCall.args[0].uri).to.equal('relay-webhooks/test');
        done();
      });
    });

    it('should throw an error if relayWebhook is null', function(done) {
      relayWebhooks.update(null, function(err) {
        expect(err.message).to.equal('relayWebhook object is required');
        expect(client.put).not.to.have.been.called;
        done();
      });
    });

    it('should throw an error if relayWebhook is missing', function(done) {
      relayWebhooks.update(function(err) {
        expect(err.message).to.equal('relayWebhook object is required');
        expect(client.put).not.to.have.been.called;
        done();
      });
    });
  });

  describe('delete Method', function() {
    it('should call client delete method with the appropriate uri', function(done) {
      relayWebhooks['delete']('test', function(err, data) {
        expect(client['delete'].firstCall.args[0].uri).to.equal('relay-webhooks/test');
        done();
      });
    });

    it('should throw an error if id is null', function(done) {
      relayWebhooks['delete'](null, function(err) {
        expect(err.message).to.equal('id is required');
        expect(client['delete']).not.to.have.been.called;
        done();
      });
    });

    it('should throw an error if id is missing', function(done) {
      relayWebhooks['delete'](function(err) {
        expect(err.message).to.equal('id is required');
        expect(client['delete']).not.to.have.been.called;
        done();
      });
    });
  });
});
