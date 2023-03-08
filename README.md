# gorala-backend

This is the corrosponding backend for [a note / calendar app written in Flutter](https://github.com/raLaaaa/gorala-client).
The backend is written in Go and uses [Echo](https://echo.labstack.com/) + [GORM](https://gorm.io/). The flutter app uses the API of this project.
Authentication is done by `JWT`.

The project also contains a password reset page as well as a basic advertisment page (the page uses google fonts for icons).
It also uses the Sendinblue API for sending emails. You can setup an environment variable called `SIBKEY` with your api key.

If you want to start the application simple use the `Dockerfile` or run `go run server.go`.

This is a pure work in progress side time project and has purely an educational purpose. 
The project is not finished how ever I decided to share it in case anyone is looking for an [Echo](https://echo.labstack.com/) API example project. 

In case you need information on how to do authentication with [Echo](https://echo.labstack.com/) and JWT, how to use GORM or how to setup an API this repository might help you.
