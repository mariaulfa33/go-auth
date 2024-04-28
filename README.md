# go-auth - Sample Project - User Management API

## Requirement
Technical Requirements:
- HTTP server with database connection
- Five API endpoints: Register, Login, List, Add User, Remove User
- Unit test for each function with 85% coverage

API Endpoints:
1. Register: An endpoint where users can register their account by providing necessary information such as username, email, and password.
2. Login: An endpoint for user authentication, using username and password as input. If authentication is successful, the endpoint should return a token for subsequent API calls.
3. List: An endpoint for retrieving a list of all users. This endpoint should require authentication using the token obtained from the Login endpoint.
4. Add User: An endpoint for adding a new user to the database, requiring the input of username, email, and password. Only authenticated users should be able to add another user.
5. Remove User: An endpoint for removing a user from the database, requiring the input of the user's ID or username. Only authenticated users should be able to remove a user.


## Deliverables:

Source code of the project with clear documentation and comments
A readme file with instructions on how to set up and run the project, including necessary dependencies and tools.

Unit tests for each function with 85% coverage and Any other necessary documentation or notes to help run or maintain the project


================================================================================

## Migration
- go install github.com/pressly/goose/v3/cmd/goose@latest
- goose postgres "DATABASE STRING CONNECTION" up


================================================================================

