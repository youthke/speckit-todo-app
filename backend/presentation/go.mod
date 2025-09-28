module todo-app/presentation

go 1.23

require (
    todo-app/domain v0.0.0
    todo-app/application v0.0.0
    todo-app/infrastructure v0.0.0
    github.com/gin-gonic/gin v1.11.0
)

replace todo-app/domain => ../domain
replace todo-app/application => ../application
replace todo-app/infrastructure => ../infrastructure

// Presentation layer: Can import all other layers for HTTP handling