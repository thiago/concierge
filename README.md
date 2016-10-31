# Concierge
:postbox: The manager all events

### Usage

#### Create file with your resource in function

> JavaScript (server.js)
```js
module.exports.test = function(req, res) {
  var messsage = { "message": "Hi, " + req.params("name") }
  return res(200, concierge.JSON, req.params("name"));
}
```

#### Get your resource in function

**GET - /{filePath}/{functionName}/{param}**
 
> Request
```bash
http://localhost:8000/server/test/concierge
```

> Response
```json
{
  "message": "Hi, concierge"
}
```
