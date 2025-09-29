module todo-app/infrastructure

go 1.23

require (
    todo-app/domain v0.0.0
    todo-app/application v0.0.0
    gorm.io/gorm v1.31.0
    gorm.io/driver/sqlite v1.6.0
)

replace todo-app/domain => ../domain
replace todo-app/application => ../application

// Infrastructure layer: Can import domain and application for implementations