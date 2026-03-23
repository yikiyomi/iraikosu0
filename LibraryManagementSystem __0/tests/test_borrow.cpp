#include <iostream>
#include "Database.hpp"
#include "BorrowManager.hpp"

void testBorrow() {
    Database db("tcp://127.0.0.1:3306", "root", "password", "library_test");
    if (!db.connect()) return;

    BorrowManager bm(db);
    // 假设存在读者ID=1，图书ISBN=...
    if (bm.borrowBook(1, "978-7-111-12345-6")) {
        std::cout << "借书测试通过" << std::endl;
    } else {
        std::cout << "借书失败" << std::endl;
    }
}

int main() {
    testBorrow();
    return 0;
}