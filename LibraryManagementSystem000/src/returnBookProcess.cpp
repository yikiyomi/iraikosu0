ReturnResult LibrarySystem::returnBookProcess(int recordId, BookCondition condition) {
    ReturnResult result;
    
    // 1. 获取借阅记录
    BorrowRecord record = BorrowRecord::fetchFromDatabase(db, recordId);
    
    // 2. 计算是否逾期
    Date currentDate = getCurrentDate();
    if (currentDate > record.dueDate) {
        int overdueDays = calculateOverdueDays(record.dueDate, currentDate);
        double fine = calculateFine(overdueDays);
        
        result.overdueDays = overdueDays;
        result.fineAmount = fine;
        
        // 更新读者罚金
        Reader reader = Reader::fetchFromDatabase(db, record.readerId);
        reader.addPenalty(fine);
        reader.updateInDatabase(db);
        
        // 添加罚款记录
        Fine fineRecord(record.readerId, fine, "逾期");
        fineRecord.saveToDatabase(db);
    }
    
    // 3. 检查书籍状况
    if (condition == BookCondition::DAMAGED) {
        double damageFine = calculateDamageFine(record.bookId);
        result.fineAmount += damageFine;
    } else if (condition == BookCondition::LOST) {
        result.fineAmount += calculateLostFine(record.bookId);
        // 更新图书总数量
        Book book = Book::fetchFromDatabase(db, record.bookId);
        book.decreaseTotalCopies();
        book.updateInDatabase(db);
    } else {
        // 正常归还，更新可借数量
        Book book = Book::fetchFromDatabase(db, record.bookId);
        book.increaseAvailableCopies();
        book.updateInDatabase(db);
    }
    
    // 4. 更新借阅记录状态
    record.returnDate = getCurrentDateTime();
    record.status = "已还";
    record.updateInDatabase(db);
    
    // 5. 更新读者借书数量
    Reader reader = Reader::fetchFromDatabase(db, record.readerId);
    reader.decreaseBorrowedCount();
    reader.updateInDatabase(db);
    
    return result;
}