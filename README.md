# Concierge
:postbox: The manager all events

### Usage

#### Create file with your resource in function

> JavaScript
```js
module.exports.test = function(req, res) {
  return res(200, concierge.JSON, "Hello World!");
}
```

#### Get your resource in function

> GET in root route passing your file, function and param name
```bash
http://localhost:8000/{file}/{functionName}/{paramName}
```
