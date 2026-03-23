#include <iostream>
#include "Database.hpp"
#include "Book.hpp"

void testAddBook() {
    Database db("tcp://127.0.0.1:3306", "root", "password", "library_test");
    if (!db.connect()) {
        std::cerr << "连接失败" << std::endl;
        return;
    }

    Book book("978-7-111-12345-6", "C++ Primer Plus", "Stephen Prata");
    book.setPublisher("机械工业出版社");
    book.setPublishDate("2020-01-01");
    book.setCategory("编程");
    book.setPrice(128.0);
    book.setTotalCopies(5);
    book.setAvailableCopies(5);

    if (book.addToDatabase(db)) {
        std::cout << "图书添加成功" << std::endl;
    } else {
        std::cout << "添加失败" << std::endl;
    }
}

int main() {
    testAddBook();
    return 0;
}