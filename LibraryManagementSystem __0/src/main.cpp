#include "LibraryUI.hpp"
#include <iostream>

int main() {
    // 数据库连接参数（请根据实际情况修改）
    Database db("tcp://127.0.0.1:3306", "root", "password", "library");

    if (!db.connect()) {
        std::cerr << "无法连接数据库，程序退出。" << std::endl;
        return 1;
    }

    LibraryUI ui(db);
    ui.mainMenu();

    return 0;
}