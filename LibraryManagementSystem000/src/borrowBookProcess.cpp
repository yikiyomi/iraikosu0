bool LibrarySystem::borrowBookProcess(int readerId, const std::string& isbn) {
    // 1. 验证读者状态
    Reader reader = Reader::fetchFromDatabase(db, readerId);
    if (!reader.canBorrowMore()) {
        std::cout << "已达到最大借书数量限制！" << std::endl;
        return false;
    }
    
    if (reader.hasOverdueBooks(db)) {
        std::cout << "有逾期图书未还！" << std::endl;
        return false;
    }
    
    // 2. 查找图书
    Book book = Book::findByISBN(db, isbn);
    if (!book.isAvailable()) {
        std::cout << "图书不可借！" << std::endl;
        return false;
    }
    
    // 3. 创建借阅记录
    BorrowRecord record;
    record.readerId = readerId;
    record.bookId = book.getBookId();
    record.borrowDate = getCurrentDateTime();
    record.dueDate = calculateDueDate(record.borrowDate, 
                                     reader.getReaderType());
    
    // 4. 更新图书状态
    book.decreaseAvailableCopies();
    book.updateInDatabase(db);
    
    // 5. 更新读者借书数量
    reader.increaseBorrowedCount();
    reader.updateInDatabase(db);
    
    // 6. 保存借阅记录
    record.saveToDatabase(db);
    
    std::cout << "借书成功！应还日期：" << record.dueDate << std::endl;
    return true;
}