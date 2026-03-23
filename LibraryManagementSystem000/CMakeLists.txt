cmake_minimum_required(VERSION 3.10)
project(LibraryManagementSystem)

set(CMAKE_CXX_STANDARD 17)

# 查找 MySQL Connector
find_package(MySQL REQUIRED)

# 添加可执行文件
add_executable(library_system
    src/main.cpp
    src/Database.cpp
    src/Book.cpp
    src/Reader.cpp
    src/BorrowManager.cpp
    src/ReportGenerator.cpp
    src/UI.cpp
)

# 链接库
target_link_libraries(library_system
    PRIVATE
    MySQL::MySQL
    ${CMAKE_THREAD_LIBS_INIT}
)