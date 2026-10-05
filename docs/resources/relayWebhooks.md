# Relay Webhooks

This library provides easy access to the [Relay Webhooks](https://www.sparkpost.com/api#/reference/relay-webhooks/) Resource.

## Methods
* **all(callback)**
  List currently extant relay webhooks.
  * `callback` - executed after task is completed. **required**
    * standard `callback(err, data)`
    * `err` - any error that occurred
    * `data` - full response from request client
* **describe(options, callback)**
  Retrieve details about a specified relay webhook by its id
  * `options.id` - the id of the relay webhook you want to describe **required**
  * `callback` - see all function
* **create(relayWebhook, callback)**
  Create a new relay webhook
  * `relayWebhook` - a relay webhook object **required**
  * `callback` - see all function
* **update(relayWebhook, callback)**
  Update an existing relay webhook
  * `relayWebhook` - a relay webhook object including its id **required**
  * `callback` - see all function
* **delete(id, callback)**
  Delete an existing relay webhook
  * `id` - the id of the relay webhook you want to delete **required**
  * `callback` - see all function

## Examples

```js
var SparkPost = require('sparkpost');
var client = new SparkPost('YOUR_API_KEY');

client.relayWebhooks.all(function(err, data) {
  if(err) {
    console.log(err);
    return;
  }

  console.log(data.body);
});

```

Check out all the examples provided [here](/examples/relayWebhooks).
