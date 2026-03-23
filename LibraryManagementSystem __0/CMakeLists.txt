cmake_minimum_required(VERSION 3.10)
project(LibraryManagementSystem)

set(CMAKE_CXX_STANDARD 17)

# 查找 MySQL Connector/C++
find_package(PkgConfig REQUIRED)
pkg_check_modules(MYSQLCPPCONN REQUIRED mysqlcppconn)

# 包含目录
include_directories(include ${MYSQLCPPCONN_INCLUDE_DIRS})

# 源文件
set(SOURCES
    src/Database.cpp
    src/Book.cpp
    src/Reader.cpp
    src/BorrowManager.cpp
    src/ReportGenerator.cpp
    src/LibraryUI.cpp
    src/main.cpp
)

# 生成可执行文件
add_executable(library_system ${SOURCES})

# 链接库
target_link_libraries(library_system ${MYSQLCPPCONN_LIBRARIES})

# 如果使用线程
find_package(Threads REQUIRED)
target_link_libraries(library_system Threads::Threads)