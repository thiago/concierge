var message = "Hi, " + request.params("name").toUpperCase();

var response = {
  "message": message
}

next(200, concierge.JSON, JSON.stringify(response));