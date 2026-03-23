#ifndef LIBRARYUI_HPP
#define LIBRARYUI_HPP

#include "Database.hpp"
#include "Book.hpp"
#include "Reader.hpp"
#include "BorrowManager.hpp"
#include "ReportGenerator.hpp"
#include <iostream>
#include <string>
#include <memory>

class LibraryUI{
    private:
    Database& db;
    BorrowManager borrowMgr;
    ReportGenerator reportGen;
    int currentAsminId;
    public:
    LibraryUI(Database& database);
    void mainMenu();
    void bookManagementMenu();
    void readerManagementMenu();
    void borrowReturnMenu();
    void reportMenu();
    void sysemMenu;

    //辅助函数
    void displayBooks(const std::vector<Book>& books);
    void displayreaders(const std::vector<Reader>& readers);
    void displayBorrowRecords(const std::vector<BorrowRecord>& records);

    //输入验证
    int getintinput(const std::string&prompt);
    std::string getstringinput(const std::string& prompt);
    std::string getdateinput(const std::string& prompt);
};

#endif