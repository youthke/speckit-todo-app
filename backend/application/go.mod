module todo-app/application

go 1.23

require (
    todo-app/domain v0.0.0
)

replace todo-app/domain => ../domain

// Application layer: Can import domain but not infrastructure or presentation