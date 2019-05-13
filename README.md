# README

Run `make` which will build the docker image and start the server on port 8080.

If your 8080 is busy please change the `HTTP_PORT` variable in the `Makefile` accordingly.

## Notes

* This API is not secured, it’s not meant to be exposed outside of the cluster. Authentication should be handled via
  another API (e.g. API Gateway).
* Pagination should be done via tokens (before, after) to avoid inconsistencies when moving forward and backwards 
  through pages.
* In production you would use a real database rather than keeping all payments in-memory. An in-memory repository could
  still be useful for testing purposes but maintaining it manually would make you write logic twice (e.g. sorting, 
  filtering...). If you're using SQL you could use an in-memory SQLite for example or if you're using Mongo then an
  in-memory MongoDB server can be easily created for your tests.
* At the moment there's no validation on the Payment model, I left it out because it requires knowledge of the domain.
  It could be done easily by hooking into the JSON Marshaling or by adding a dedicated service to validate payments that
  could be injected via the container like I'm doing with the repository.

## Tests

The unit tests are executed when the Docker image is built right before compiling the go binary. 
Once the program is compiled it is copied in a `scratch` image without the rest.
The final image weighs 11.2MB (production ready). 

If you want to test the API manually I've added a `Postman_collection.json` inside the `swagger` folder
that you can import in your Postman to make things easy. Just remember to create a `gopaytest-api` variable in
your Postman and set it to wherever your server is listening to (e.g. `http://localhost:8080`).
